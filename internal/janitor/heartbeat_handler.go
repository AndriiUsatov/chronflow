package janitor

import (
	"fmt"
	"net/http"

	config "github.com/AndriiUsatov/chronflow/internal/config/janitorcfg"
	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	heartbeatEnpoint = "/heartbeat"
	metricsEndpoint  = "/metrics"
)

type heartbeatServer struct {
	addr    string
	repo    db.TaskRepository
	handler *gin.Engine
	Metrics *JanitorMetrics
}

type JanitorMetrics struct {
	TaskRecovered prometheus.Counter
}

func (server heartbeatServer) ListenAndServe() error {
	ser := http.Server{
		Addr:    server.addr,
		Handler: server.handler,
	}

	return ser.ListenAndServe()
}

func getMetrics() *JanitorMetrics {
	return &JanitorMetrics{
		TaskRecovered: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "chronflow_janitor_task_recovered",
				Help: "Total count of tasks recovered by Janitor service",
			},
		),
	}
}

func NewHeartBeatHandler(cfg config.JanitorConfig, repo db.TaskRepository) heartbeatServer {
	r := gin.Default()

	res := heartbeatServer{
		addr:    fmt.Sprintf(":%s", cfg.JanitorHeartbeatPort),
		repo:    repo,
		handler: r,
		Metrics: getMetrics(),
	}

	r.GET(heartbeatEnpoint, func(ctx *gin.Context) {
		err := repo.Ping(ctx)
		if err == nil {
			ctx.Status(http.StatusOK)
		}
	})

	r.GET(metricsEndpoint, gin.WrapH(promhttp.Handler()))

	return res
}
