package scheduler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsEndpoint = "/metrics"
)

type metricsServer struct {
	engine  *gin.Engine
	port    int
	Metrics *SchedulerMetrics
}

type SchedulerMetrics struct {
	taskScheduler prometheus.Counter
}

func (srv metricsServer) ListenAndServe() {
	srv.engine.Run()
}

func getMetrics() *SchedulerMetrics {
	return &SchedulerMetrics{
		taskScheduler: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "chronflow_scheduler_task_scheduled",
				Help: "Total count of scheduled tasks by Scheduler service",
			},
		),
	}
}

func NewMetricsServer(port int) *metricsServer {
	r := gin.Default()

	result := &metricsServer{
		engine:  r,
		port:    port,
		Metrics: getMetrics(),
	}

	r.GET(metricsEndpoint, gin.WrapH(promhttp.Handler()))

	return result
}
