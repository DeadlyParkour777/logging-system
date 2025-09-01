package main

import (
	"context"
	"flag"
	"log"
	"test-client/internal/app"
	"test-client/internal/client"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/logs"
	"github.com/DeadlyParkour777/logging-system/pkg/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddr = "localhost:8081"
const otelCollectorAddr = "localhost:4317"

func main() {
	mode := flag.String("mode", "single", "Test mode: 'single' or 'loop'")
	level := flag.String("level", "INFO", "Log level for single mode (INFO, WARN, ERROR)")
	message := flag.String("message", "This is a test message.", "Log message for single mode")
	serviceName := flag.String("service", "test-client", "Service name to report in logs")
	interval := flag.Duration("interval", 5*time.Second, "Interval for loop mode (e.g., 5s, 1m)")
	flag.Parse()

	tp, err := telemetry.InitTracer(*serviceName, otelCollectorAddr)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	conn, err := grpc.NewClient(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("FATAL: did not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient := client.NewGRPCClient(logs.NewLogServiceClient(conn))
	app := app.New(grpcClient)

	app.Run(context.Background(), *mode, *level, *message, *serviceName, *interval)
}
