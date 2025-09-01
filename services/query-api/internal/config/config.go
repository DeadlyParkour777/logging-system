package config

import (
	"log"
	"os"
	"time"

	"go.yaml.in/yaml/v2"
)

type Config struct {
	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`

	ClickHouse struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"clickhouse"`

	Redis struct {
		Address         string        `yaml:"address"`
		Password        string        `yaml:"password"`
		DB              int           `yaml:"db"`
		CacheTTlSeconds time.Duration `yaml:"cache_ttl_seconds"`
	} `yaml:"redis"`
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
