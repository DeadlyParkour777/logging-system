package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka struct {
		Brokers     []string `yaml:"brokers"`
		Topic       string   `yaml:"topic"`
		AlertsTopic string   `yaml:"alerts_topic"`
		GroupID     string   `yaml:"group_id"`
	} `yaml:"kafka"`
	ClickHouse struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"clickhouse"`
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

	return cfg
}
