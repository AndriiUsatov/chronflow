package janitor

import (
	"fmt"
	"net/http"

	config "github.com/AndriiUsatov/chronflow/internal/config/janitorcfg"
	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/gin-gonic/gin"
)

type heartbeatServer struct {
	addr    string
	repo    db.TaskRepository
	handler *gin.Engine
}

func (server heartbeatServer) ListenAndServe() error {
	ser := http.Server{
		Addr:    server.addr,
		Handler: server.handler,
	}

	return ser.ListenAndServe()
}

func NewHeartBeatHandler(cfg config.JanitorConfig, repo db.TaskRepository) heartbeatServer {
	r := gin.Default()

	res := heartbeatServer{
		addr:    fmt.Sprintf(":%s", cfg.JanitorHeartbeatPort),
		repo:    repo,
		handler: r,
	}

	r.GET("/heartbeat", func(ctx *gin.Context) {
		err := repo.Ping(ctx)
		if err == nil {
			ctx.Status(http.StatusOK)
		}
	})

	return res
}
