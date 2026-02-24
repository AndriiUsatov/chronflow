package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/AndriiUsatov/chronflow/internal/config/janitorcfg"
	"github.com/AndriiUsatov/chronflow/internal/db/postgres"
	"github.com/AndriiUsatov/chronflow/internal/janitor"
)

func main() {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Configuration
	cfg, err := config.LoadApiConfig("./")
	if err != nil {
		panic(err)
	}

	// Task repository
	taskRepo, err := postgres.NewPGRepositoryWithRetry(cfg.PGURI, cfg.PGPort, cfg.PGUser, cfg.PGPwd, cfg.PGTaskDB, cfg.PGSSLMode, cfg.PGTaskSchema, cfg.PGTaskTable, ctx, 5*time.Minute)
	if err != nil {
		panic(err)
	}
	defer taskRepo.Db.Close()

	// Hearbeat
	go janitor.NewHeartBeatHandler(cfg, taskRepo).ListenAndServe()

	// Janitor
	jntr := postgres.NewPostgresJanitor(taskRepo)
	go jntr.Run(ctx)

	<-ctx.Done()
}
