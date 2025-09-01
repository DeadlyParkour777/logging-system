package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
)

type Store interface {
	SearchLogsInDB(ctx context.Context, params map[string]string) ([]models.LogEntry, error)
	GetLogsFromCache(ctx context.Context, key string) ([]models.LogEntry, error)
	SetLogsInCache(ctx context.Context, key string, logs []models.LogEntry, ttl time.Duration) error
}

type QueryService struct {
	store    Store
	cacheTTL time.Duration
}

func NewQueryService(store Store, cacheTTL time.Duration) *QueryService {
	return &QueryService{store: store, cacheTTL: cacheTTL}
}

func (s *QueryService) SearchLogs(ctx context.Context, params map[string]string) ([]models.LogEntry, error) {
	cacheKey := generateCacheKey(params)
	log.Printf("Searching with cache key: %s", cacheKey)

	cachedLogs, err := s.store.GetLogsFromCache(ctx, cacheKey)
	if err != nil {
		log.Printf("WARNING: Redis Get failed: %v. Fetching from DB.", err)
	}
	if cachedLogs != nil {
		return cachedLogs, nil
	}

	dbLogs, err := s.store.SearchLogsInDB(ctx, params)
	if err != nil {
		return nil, err
	}

	if err := s.store.SetLogsInCache(ctx, cacheKey, dbLogs, s.cacheTTL); err != nil {
		log.Printf("WARNING: Redis Set failed: %v", err)
	}

	return dbLogs, nil
}

func generateCacheKey(params map[string]string) string {
	rawKey := fmt.Sprintf("s:%s_l:%s_q:%s", params["service_name"], params["level"], params["search"])
	return fmt.Sprintf("logs:%x", sha1.Sum([]byte(rawKey)))
}
