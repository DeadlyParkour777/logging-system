package config

import (
	"log"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

type Config struct {
	GRPC struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
	} `yaml:"kafka"`
}

func NewConfig(configPath string) *Config {
	rawYAML, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(rawYAML, cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return cfg
}
