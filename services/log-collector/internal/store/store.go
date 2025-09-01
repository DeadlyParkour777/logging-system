package store

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type KafkaStore struct {
	writer *kafka.Writer
}

func NewKafkaStore(brokers []string, topic string) *KafkaStore {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},

		RequiredAcks: kafka.RequireOne,

		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	}

	return &KafkaStore{writer: writer}
}

func (s *KafkaStore) SaveLog(ctx context.Context, entry models.LogEntry) error {
	tracer := otel.Tracer("log-collector-store")
	ctx, span := tracer.Start(ctx, "SaveLogToKafka")
	defer span.End()

	msgBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("ERROR: marshal log entry: %v", err)
		span.RecordError(err)
		return err
	}
	span.AddEvent("Log entry marshaled to JSON")

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
		attribute.String("kafka.topic", s.writer.Topic),
		attribute.Int("kafka.message_size", len(msgBytes)),
	)

	err = s.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("ERROR: could not write message to kafka: %v", err)
		span.RecordError(err)
		return err
	}

	log.Printf("Log sent to kafka topic '%s'", s.writer.Topic)
	return nil
}

func (s *KafkaStore) Close() error {
	log.Println("Closing kafka writer...")
	return s.writer.Close()
}
