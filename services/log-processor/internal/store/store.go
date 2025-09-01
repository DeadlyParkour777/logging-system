package store

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type ClickhouseStore struct {
	conn clickhouse.Conn
}

type AlertProducer struct {
	writer *kafka.Writer
}

func NewClickhouseStore(addr, db, user, pass string) (*ClickhouseStore, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
	})

	if err != nil {
		return nil, err
	}
	return &ClickhouseStore{conn: conn}, nil
}

func (s *ClickhouseStore) SaveLog(ctx context.Context, entry models.LogEntry) error {
	tracer := otel.Tracer("log-processor-store")
	ctx, span := tracer.Start(ctx, "SaveLogInClickHouse")
	defer span.End()

	query := `INSERT INTO logs (timestamp, service_name, level, message, metadata) VALUES (?, ?, ?, ?, ?)`
	err := s.conn.Exec(ctx, query, entry.Timestamp, entry.ServiceName, entry.Level, entry.Message, entry.Metadata)
	if err != nil {
		log.Printf("ERROR: failed to save log to Clickhouse: %v", err)
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *ClickhouseStore) Close() error {
	return s.conn.Close()
}

func NewAlertProducer(brokers []string, topic string) *AlertProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &AlertProducer{writer: writer}
}

func (p *AlertProducer) PublishAlert(ctx context.Context, entry models.LogEntry) error {
	tracer := otel.Tracer("log-processor-store")
	ctx, span := tracer.Start(ctx, "PublishAlertToKafka")
	defer span.End()

	msgBytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.MapCarrier(headers))

	kafkaHeaders := make([]kafka.Header, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: k, Value: []byte(v)})
	}

	msg := kafka.Message{
		Value:   msgBytes,
		Headers: kafkaHeaders,
	}

	span.SetAttributes(
		attribute.String("kafka.topic", p.writer.Topic),
	)

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("ERROR: could not publish alert to kafka: %v", err)
		return err
	}

	log.Printf("Alert for service '%s' published to topic '%s'", entry.ServiceName, p.writer.Topic)
	return nil
}

func (p *AlertProducer) Close() error {
	return p.writer.Close()
}
