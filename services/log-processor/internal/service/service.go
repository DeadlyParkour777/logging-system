package service

import (
	"context"
	"log"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"go.opentelemetry.io/otel"
)

type LogStore interface {
	SaveLog(ctx context.Context, entry models.LogEntry) error
}

type AlertPublisher interface {
	PublishAlert(ctx context.Context, entry models.LogEntry) error
}

type LogProcessorService struct {
	store   LogStore
	alerter AlertPublisher
}

func NewLogProcessorService(store LogStore, alerter AlertPublisher) *LogProcessorService {
	return &LogProcessorService{store: store, alerter: alerter}
}

func (s *LogProcessorService) ProcessLog(ctx context.Context, entry models.LogEntry) error {
	tracer := otel.Tracer("log-processor-service")
	ctx, span := tracer.Start(ctx, "ProcessLog")
	defer span.End()

	if err := s.store.SaveLog(ctx, entry); err != nil {
		span.RecordError(err)
		return err
	}

	if entry.Level == "ERROR" || entry.Level == "CRITICAL" {
		span.AddEvent("Detected ERROR log, preparing to publish alert")
		if err := s.alerter.PublishAlert(ctx, entry); err != nil {
			log.Printf("WARNING: failed to publish alert. Error: %v", err)
			span.RecordError(err)
		}
	}

	return nil
}
