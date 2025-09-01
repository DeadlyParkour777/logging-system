package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type LogSender interface {
	Send(ctx context.Context, req *logs.SendLogRequest) error
}

type GRPCClient struct {
	client logs.LogServiceClient
}

func NewGRPCClient(client logs.LogServiceClient) *GRPCClient {
	return &GRPCClient{client: client}
}

func (c *GRPCClient) Send(ctx context.Context, req *logs.SendLogRequest) error {
	tracer := otel.Tracer("test-client")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("Send %s Log", req.Level))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	log.Printf("ending log: Level=%s, Service=%s, Message='%s'", req.Level, req.ServiceName, req.Message)

	span.SetAttributes(
		attribute.String("request.service_name", req.ServiceName),
		attribute.String("request.level", req.Level),
	)

	res, err := c.client.SendLog(ctx, req)
	if err != nil {
		span.RecordError(err)
		return err
	}

	span.SetAttributes(attribute.Bool("response.ok", res.GetOk()))
	log.Printf("Response received: OK = %t", res.GetOk())
	return nil
}
