package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/api"
	config "github.com/AndriiUsatov/chronflow/internal/config/apicfg"
	"github.com/AndriiUsatov/chronflow/internal/db/postgres"
)

func init() {
	time.Local = time.UTC
}

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

	// Task REST API server
	go api.GetTaskRestServer(cfg, taskRepo).ListenAndServe(true)

	// gRPC status update server
	statusServer, err := api.NewTaskStatusServer(cfg.GrpcTaskUpdateServerTransportProtocol, fmt.Sprintf(":%d", cfg.GrpcTaskUpdateServerPort), taskRepo)
	if err != nil {
		panic(err)
	}
	defer statusServer.Close()
	go statusServer.ListenAndServe(true)

	<-ctx.Done()
}
