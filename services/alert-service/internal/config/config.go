package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v2"
)

type Config struct {
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`
	Telegram struct {
		BotToken string
		ChatID   string
	}
}

func NewConfig(path string) *Config {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	cfg.Telegram.BotToken = os.Getenv("TELEGRAM_TOKEN")
	cfg.Telegram.ChatID = os.Getenv("TELEGRAM_CHAT_ID")

	return cfg
}
