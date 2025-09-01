package service

import (
	"context"
	"errors"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"go.opentelemetry.io/otel"
)

type LogStore interface {
	SaveLog(ctx context.Context, entry models.LogEntry) error
}

type CollectorService struct {
	store LogStore
}

func NewCollectorService(store LogStore) *CollectorService {
	return &CollectorService{store: store}
}

func (s *CollectorService) ProcessSaveLog(ctx context.Context, entry models.LogEntry) error {
	tracer := otel.Tracer("log-collector-service")
	ctx, span := tracer.Start(ctx, "ProcessSaveLog")
	defer span.End()

	if entry.ServiceName == "" {
		err := errors.New("service_name is required")
		span.RecordError(err)
		return err
	}

	if entry.Message == "" {
		err := errors.New("message is required")
		span.RecordError(err)
		return err
	}

	entry.Timestamp = time.Now().UTC()

	span.AddEvent("Log enriched with timestamp")

	return s.store.SaveLog(ctx, entry)
}
