package handler

import (
	"context"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CollectorService interface {
	ProcessSaveLog(ctx context.Context, entry models.LogEntry) error
}

type GrpcHandler struct {
	service CollectorService
	logs.UnimplementedLogServiceServer
}

func NewGrpcHandler(service CollectorService) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) SendLog(ctx context.Context, req *logs.SendLogRequest) (*logs.SendLogResponse, error) {
	tracer := otel.Tracer("log-collector-handler")
	ctx, span := tracer.Start(ctx, "SendLog")
	defer span.End()

	entry := models.LogEntry{
		ServiceName: req.GetServiceName(),
		Level:       req.GetLevel(),
		Message:     req.GetMessage(),
		Metadata:    req.GetMetadata(),
	}

	span.SetAttributes(
		attribute.String("log.service_name", entry.ServiceName),
		attribute.String("log.level", entry.Level),
	)

	if err := h.service.ProcessSaveLog(ctx, entry); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to process log")
	}

	return &logs.SendLogResponse{Ok: true}, nil
}
