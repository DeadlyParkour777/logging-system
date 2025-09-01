package models

import "time"

type LogEntry struct {
	ServiceName string            `json:"service_name"`
	Level       string            `json:"level"`
	Message     string            `json:"message"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
