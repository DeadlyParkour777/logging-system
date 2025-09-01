package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/DeadlyParkour777/logging-system/services/query-api/docs"
	"github.com/DeadlyParkour777/logging-system/services/query-api/internal/config"
	"github.com/DeadlyParkour777/logging-system/services/query-api/internal/handler"
	"github.com/DeadlyParkour777/logging-system/services/query-api/internal/service"
	"github.com/DeadlyParkour777/logging-system/services/query-api/internal/store"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

const cfgPath = "/app/config.yaml"

// @title Logging System - Query API
// @version 1.0
// @description API для поиска и фильтрации логов, сохраненных в системе.
// @host localhost:8082
// @BasePath /
func main() {
	cfg := config.NewConfig(cfgPath)

	chConn, err := connectToClickHouse(cfg)
	if err != nil {
		log.Fatalf("failed to connect to clickhouse: %v", err)
	}
	redisClient := connectToRedis(cfg)

	store := store.NewStore(chConn, redisClient)
	queryService := service.NewQueryService(store, cfg.Redis.CacheTTlSeconds*time.Second)
	httpHandler := handler.NewHttpHandler(queryService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/logs", httpHandler.SearchLogs)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}

func connectToClickHouse(cfg *config.Config) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouse.Address},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.Database,
			Username: cfg.ClickHouse.Username,
			Password: cfg.ClickHouse.Password,
		},
	})
	if err != nil {
		return nil, err
	}
	return conn, conn.Ping(context.Background())
}

func connectToRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	return rdb
}
