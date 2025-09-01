package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/segmentio/kafka-go"
)

type AlertService interface {
	ProcessAlert(ctx context.Context, entry models.LogEntry) error
}

type KafkaHandler struct {
	reader  *kafka.Reader
	service AlertService
}

func NewKafkaHandler(brokers []string, topic, groupID string, service AlertService) *KafkaHandler {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &KafkaHandler{reader: reader, service: service}
}

func (h *KafkaHandler) Run(ctx context.Context) {
	defer h.reader.Close()
	for {
		msg, err := h.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("can't fetch message: %v", err)
		}

		var entry models.LogEntry
		if err := json.Unmarshal(msg.Value, &entry); err != nil {
			log.Printf("can't unmarshall message: %v", err)
			h.reader.CommitMessages(ctx, msg)
			continue
		}

		log.Printf("Alert message received for service: %s", entry.ServiceName)
		if err := h.service.ProcessAlert(ctx, entry); err != nil {
			log.Printf("failed to process alert: %v", err)
		}

		h.reader.CommitMessages(ctx, msg)
	}
}
