package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test-client/internal/client"
	"test-client/internal/generator"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
	"go.opentelemetry.io/otel"
)

type App struct {
	logSender client.LogSender
}

func New(logSender client.LogSender) *App {
	return &App{logSender: logSender}
}

func (a *App) Run(ctx context.Context, mode, level, message, serviceName string, interval time.Duration) {
	tracer := otel.Tracer("test-client-app")
	ctx, span := tracer.Start(ctx, "RunTestClientMode")
	defer span.End()

	switch mode {
	case "single":
		log.Printf("Running in SINGLE mode. Service: %s, Level: %s", serviceName, level)
		req := &logs.SendLogRequest{
			ServiceName: serviceName,
			Level:       level,
			Message:     message,
			Metadata:    map[string]string{"mode": "single"},
		}
		if err := a.logSender.Send(ctx, req); err != nil {
			log.Fatalf("Failed to send single log: %v", err)
			span.RecordError(err)
		}
	case "loop":
		log.Printf("Running in LOOP mode. Interval: %v. Press Ctrl+C to stop.", interval)
		a.runLoopMode(ctx, serviceName, interval)
	default:
		log.Fatalf("Unknown mode: %s. Use 'single' or 'loop'.", mode)
	}
}

func (a *App) runLoopMode(ctx context.Context, serviceName string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			req := generator.GenerateRandomLog(serviceName)
			if err := a.logSender.Send(ctx, req); err != nil {
				log.Printf("ERROR sending log: %v", err)
			}
		case <-stop:
			log.Println("Shutting down loop mode...")
			return
		}
	}
}
