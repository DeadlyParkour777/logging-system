package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
	"github.com/DeadlyParkour777/logging-system/pkg/telemetry"
	"github.com/DeadlyParkour777/logging-system/services/log-collector/internal/config"
	"github.com/DeadlyParkour777/logging-system/services/log-collector/internal/handler"
	"github.com/DeadlyParkour777/logging-system/services/log-collector/internal/service"
	"github.com/DeadlyParkour777/logging-system/services/log-collector/internal/store"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const cfgPath = "/app/config.yaml"

func main() {
	cfg := config.NewConfig(cfgPath)

	tp, err := telemetry.InitTracer("log-collector", "otel-collector:4317")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	logStore := store.NewKafkaStore(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer logStore.Close()

	collectorServie := service.NewCollectorService(logStore)
	grpcHandler := handler.NewGrpcHandler(collectorServie)

	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	logs.RegisterLogServiceServer(grpcServer, grpcHandler)

	log.Printf("grpc server started at %v", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
