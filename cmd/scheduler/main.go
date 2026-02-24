package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/AndriiUsatov/chronflow/internal/config/schedulercfg"
	"github.com/AndriiUsatov/chronflow/internal/db/postgres"
	"github.com/AndriiUsatov/chronflow/internal/queue/nats"
	"github.com/AndriiUsatov/chronflow/internal/scheduler"
)

func init() {
	time.Local = time.UTC
}

func main() {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Configuration
	cfg, err := config.LoadSchedulerConfig("./")
	if err != nil {
		panic(err)
	}

	// Repository
	taskRepo, err := postgres.NewPGRepositoryWithRetry(cfg.PGURI, cfg.PGPort, cfg.PGUser, cfg.PGPwd, cfg.PGTaskDB, cfg.PGSSLMode, cfg.PGTaskSchema, cfg.PGTaskTable, ctx, time.Minute*5)
	if err != nil {
		panic(err)
	}
	defer taskRepo.Db.Close()

	// Notifier (New tasks arrived)
	notifier, err := postgres.NewPGListener(ctx, cfg.PGURI, cfg.PGPort, cfg.PGUser, cfg.PGPwd, cfg.PGTaskDB, cfg.PGSSLMode, cfg.PGNotificationEvent, time.Minute*5)
	if err != nil {
		panic(err)
	}
	defer notifier.Close()

	// Task Queue (NATS queue) for sending Tasks from `Scheduler` to `Worker`
	queue, err := nats.NewQueue(ctx, cfg.NatsURL, cfg.NatsTaskStream, cfg.NatsTaskToProcessSubject, cfg.NatsTaskToProcessDurable, time.Minute*5)
	if err != nil {
		panic(err)
	}
	defer queue.Close()

	// Collects Tasks from Db and send it to `Worker`
	tsch := scheduler.TaskScheduler{
		Repo:     taskRepo,
		Notifier: notifier,
		Queue:    queue,
	}

	err = tsch.Run(ctx)
	if err != nil {
		panic(err)
	}

}
