package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	ch    clickhouse.Conn
	redis *redis.Client
}

func NewStore(ch clickhouse.Conn, rds *redis.Client) *Store {
	return &Store{ch: ch, redis: rds}
}

func (s *Store) SearchLogsInDB(ctx context.Context, params map[string]string) ([]models.LogEntry, error) {
	var conditions []string
	var args []any

	if val, ok := params["service_name"]; ok && val != "" {
		conditions = append(conditions, "service_name = ?")
		args = append(args, val)
	}

	if val, ok := params["level"]; ok && val != "" {
		conditions = append(conditions, "level = ?")
		args = append(args, val)
	}

	if val, ok := params["search"]; ok && val != "" {
		conditions = append(conditions, "message LIKE ?")
		args = append(args, "%"+val+"%")
	}

	query := "SELECT timestamp, service_name, level, message, metadata FROM logs"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY timestamp DESC LIMIT 100"

	rows, err := s.ch.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("clickhouse query failed: %w", err)
	}
	defer rows.Close()

	var results []models.LogEntry
	for rows.Next() {
		var entry models.LogEntry
		if err := rows.Scan(&entry.Timestamp, &entry.ServiceName, &entry.Level, &entry.Message, &entry.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, entry)
	}

	return results, nil
}

func (s *Store) GetLogsFromCache(ctx context.Context, key string) ([]models.LogEntry, error) {
	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var results []models.LogEntry
	if err := json.Unmarshal([]byte(val), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Store) SetLogsInCache(ctx context.Context, key string, logs []models.LogEntry, ttl time.Duration) error {
	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, ttl).Err()
}
