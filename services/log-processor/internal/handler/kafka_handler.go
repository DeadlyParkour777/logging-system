package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type LogProcessorService interface {
	ProcessLog(ctx context.Context, entry models.LogEntry) error
}

type KafkaHandler struct {
	reader  *kafka.Reader
	service LogProcessorService
}

func NewKafkaHandler(brokers []string, topic, groupID string, service LogProcessorService) *KafkaHandler {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 1e6,
		MaxWait:  1 * time.Second,
	})

	return &KafkaHandler{
		reader:  reader,
		service: service,
	}
}

func (h *KafkaHandler) Run(ctx context.Context) {
	defer h.reader.Close()

	for {
		msg, err := h.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("ERROR: can't fetch message: %v", err)
			continue
		}

		headers := make(map[string]string)
		for _, h := range msg.Headers {
			headers[h.Key] = string(h.Value)
		}
		propagator := otel.GetTextMapPropagator()
		msgCtx := propagator.Extract(context.Background(), propagation.MapCarrier(headers))

		tracer := otel.Tracer("log-processor-handler")
		ctx, span := tracer.Start(msgCtx, "ProcessKafkaMessage")

		var entry models.LogEntry
		if err := json.Unmarshal(msg.Value, &entry); err != nil {
			log.Printf("ERROR: can't unmarshall message value: %v", err)
			h.reader.CommitMessages(ctx, msg)
			continue
		}

		log.Printf("Message received: service=%s", entry.ServiceName)

		if err := h.service.ProcessLog(ctx, entry); err != nil {
			log.Printf("ERROR: failed to process log: %v", err)
		}

		span.End()

		if err := h.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("ERROR: failed to commit message: %v", err)
		}
	}
}
