package main

import (
	"context"
	"log"

	"github.com/DeadlyParkour777/logging-system/log-processor/internal/config"
	"github.com/DeadlyParkour777/logging-system/log-processor/internal/handler"
	"github.com/DeadlyParkour777/logging-system/log-processor/internal/service"
	"github.com/DeadlyParkour777/logging-system/log-processor/internal/store"
	"github.com/DeadlyParkour777/logging-system/pkg/telemetry"
)

func main() {
	tp, err := telemetry.InitTracer("log-processor", "otel-collector:4317")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	cfg := config.NewConfig("/app/config.yaml")

	logStore, err := store.NewClickhouseStore(
		cfg.ClickHouse.Address,
		cfg.ClickHouse.Database,
		cfg.ClickHouse.Username,
		cfg.ClickHouse.Password,
	)
	if err != nil {
		log.Fatalf("failed to connect to clickhouse: %v", err)
	}
	defer logStore.Close()

	alertProducer := store.NewAlertProducer(cfg.Kafka.Brokers, cfg.Kafka.AlertsTopic)

	logProcessorService := service.NewLogProcessorService(logStore, alertProducer)
	kafkaHandler := handler.NewKafkaHandler(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
		logProcessorService,
	)
	kafkaHandler.Run(context.Background())
}
