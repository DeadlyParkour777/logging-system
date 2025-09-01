package main

import (
	"context"

	"github.com/DeadlyParkour777/logging-system/services/alert-service/internal/config"
	"github.com/DeadlyParkour777/logging-system/services/alert-service/internal/handler"
	"github.com/DeadlyParkour777/logging-system/services/alert-service/internal/service"
	"github.com/DeadlyParkour777/logging-system/services/alert-service/internal/store"
)

const configPath = "/app/config.yaml"

func main() {
	cfg := config.NewConfig(configPath)

	notifier := store.NewTelegramNotifier(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	alertService := service.NewAlertService(notifier)
	kafkaHandler := handler.NewKafkaHandler(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, alertService)
	kafkaHandler.Run(context.Background())
}
