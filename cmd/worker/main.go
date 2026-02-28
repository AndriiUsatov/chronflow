package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/AndriiUsatov/chronflow/internal/config/workercfg"
	"github.com/AndriiUsatov/chronflow/internal/queue/nats"
	"github.com/AndriiUsatov/chronflow/internal/worker"
)

func init() {
	time.Local = time.UTC
}

func main() {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Configuration
	cfg, err := config.LoadWorkerConfig("./")
	if err != nil {
		panic(err)
	}

	// Task Queue (NATS queue) for receiving Tasks from `Scheduler`
	queue, err := nats.NewQueue(ctx, cfg.NatsURL, cfg.NatsTaskStream, cfg.NatsTaskToProcessSubject, cfg.NatsTaskToProcessDurable, 5*time.Minute)
	if err != nil {
		panic(err)
	}
	defer queue.Close()

	// Metrics server
	metricsServer := worker.NewMetricsServer(cfg.WorkerPort)
	go metricsServer.ListenAndServe()

	// Worker
	wrk, err := worker.NewTaskWorker(ctx, queue, fmt.Sprintf("%s:%d", cfg.GrpcTaskUpdateServerURL, cfg.GrpcTaskUpdateServerPort), 5*time.Minute, metricsServer.Metrics)
	if err != nil {
		panic(err)
	}
	wrk.Run(ctx)
}
