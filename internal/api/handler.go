package api

// @title           Chronflow API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @host      localhost:80
// @BasePath  /api/v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/AndriiUsatov/chronflow/docs"
	config "github.com/AndriiUsatov/chronflow/internal/config/apicfg"
	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	home                = "/"
	taskEndpoint        = "/api/v1/task"
	taskByIDEndpoint    = "/api/v1/task/:id"
	taskSwaggerEndpoint = "/swagger/*any"
	swaggerUI           = "/swagger/index.html"
	metricsEndpoint     = "/metrics"
)

const (
	statusLabel   = "status"
	statusFail    = "fail"
	statusSuccess = "success"
)

type taskAPIHandler struct {
	addr    string
	handler http.Handler
	repo    db.TaskRepository
	metrics apiMetrics
}

func (api taskAPIHandler) ListenAndServe(isPanic bool) error {
	server := http.Server{
		Addr:    api.addr,
		Handler: api.handler,
	}
	err := server.ListenAndServe()
	if isPanic && err != nil {
		panic(err)
	}
	return err
}

type taskInput struct {
	URL     string      `json:"url" binding:"required"`
	Method  string      `json:"method" binding:"required"`
	RunAt   time.Time   `json:"runAt" binding:"required"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

type taskCreateResponse struct {
	TaskId string
}

type APIError struct {
	Error string `json:"error"`
}

type taskView struct {
	ID uuid.UUID

	URL     string
	Method  string
	Headers model.JSONHeader
	Body    string

	Status       model.TaskStatus
	RunAt        time.Time
	RetryCount   int
	Created      time.Time
	Updated      time.Time
	ErrorMessage string
}

// GetTaskByID godoc
// @Summary      Show an task
// @Description  Get task by UUID
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Task UUID"
// @Success      200  {object}  taskView
// @Failure      400  {object}  APIError
// @Failure      404  {object}  APIError
// @Failure      500  {object}  APIError
// @Router       /task/{id} [get]
func (api *taskAPIHandler) getTaskByIDHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{err.Error()})
	}

	task, err := api.repo.GetTaskByUUID(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, taskView{
		ID:           task.ID,
		Method:       task.Method,
		Headers:      task.Headers,
		Body:         string(task.Body),
		Status:       task.Status,
		RunAt:        task.RunAt,
		RetryCount:   task.RetryCount,
		Created:      task.Created,
		Updated:      task.Updated,
		ErrorMessage: task.ErrorMessage,
	})
}

// AddTask godoc
// @Summary      Add task
// @Description  Creates task to be processed
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        task  body  taskInput  true  "Task details"
// @Success      200  {object}  taskCreateResponse
// @Failure      400  {object}  APIError
// @Failure      404  {object}  APIError
// @Failure      500  {object}  APIError
// @Router       /task/ [post]
func (api *taskAPIHandler) addTaskHandler(ctx *gin.Context) {
	var tsk taskInput
	if err := ctx.ShouldBindJSON(&tsk); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{err.Error()})
		return
	}
	id := uuid.New()

	err := api.repo.CreateTask(
		ctx.Request.Context(),
		model.Task{
			ID:      id,
			URL:     tsk.URL,
			Method:  strings.ToUpper(tsk.Method),
			RunAt:   tsk.RunAt,
			Headers: model.JSONHeader(tsk.Headers),
			Body:    []byte(tsk.Body),
		},
	)

	if err != nil {
		api.metrics.taskProcessed.WithLabelValues(statusFail).Inc()
		ctx.JSON(http.StatusInternalServerError, APIError{err.Error()})
		return
	}

	api.metrics.taskProcessed.WithLabelValues(statusSuccess).Inc()

	ctx.JSON(http.StatusOK, taskCreateResponse{id.String()})
}

type apiMetrics struct {
	taskProcessed *prometheus.CounterVec
}

func GetTaskRestServer(cfg config.ApiConfig, taskRepo db.TaskRepository) *taskAPIHandler {
	router := gin.Default()

	result := &taskAPIHandler{
		addr:    fmt.Sprintf(":%s", cfg.TaskAPIPort),
		handler: router,
		repo:    taskRepo,
		metrics: apiMetrics{
			taskProcessed: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "chronflow_api_task_processed",
					Help: "Total count of processed task by API service",
				},
				[]string{statusLabel},
			),
		},
	}

	router.GET(taskSwaggerEndpoint, ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET(taskByIDEndpoint, result.getTaskByIDHandler)
	router.POST(taskEndpoint, result.addTaskHandler)
	router.GET(home, func(ctx *gin.Context) {
		ctx.Redirect(http.StatusSeeOther, swaggerUI)
	})

	router.GET(metricsEndpoint, gin.WrapH(promhttp.Handler()))

	return result

}
