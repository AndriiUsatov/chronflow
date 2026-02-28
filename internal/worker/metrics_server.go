package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsEndpoint = "/metrics"
	statusLabel     = "status"
	statusFail      = "fail"
	statusSuccess   = "success"
)

type metricsServer struct {
	engine  *gin.Engine
	port    int
	Metrics *WorkerMetrics
}

type WorkerMetrics struct {
	taskProcessed *prometheus.CounterVec
}

func (srv metricsServer) ListenAndServe() {
	srv.engine.Run()
}

func getMetrics() *WorkerMetrics {
	return &WorkerMetrics{
		taskProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "chronflow_worker_task_processed",
				Help: "Total number of processed tasks by Worker service",
			},
			[]string{statusLabel},
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
