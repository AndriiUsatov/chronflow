package api

import (
	"fmt"
	"net/http"
	"time"

	config "github.com/AndriiUsatov/chronflow/internal/config/apicfg"
	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	taskEndpoint     = "/api/v1/task"
	taskByIDEndpoint = "/api/v1/task/:id"
)

type taskAPIHandler struct {
	addr    string
	handler http.Handler
	repo    db.TaskRepository
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
	URL     string      `form:"url" binding:"required"`
	Method  string      `form:"method" binding:"required"`
	RunAt   time.Time   `form:"runAt" binding:"required" time_format:"2006-01-02 15:04:05"`
	Headers http.Header `form:"headers"`
	Body    []byte      `form:"body"`
}

func (api taskAPIHandler) getTaskByIDHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	task, err := api.repo.GetTaskByUUID(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (api taskAPIHandler) addTaskHandler(ctx *gin.Context) {
	var tsk taskInput
	if err := ctx.ShouldBind(&tsk); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New()

	err := api.repo.CreateTask(
		ctx.Request.Context(),
		model.Task{
			ID:      id,
			URL:     tsk.URL,
			Method:  tsk.Method,
			RunAt:   tsk.RunAt,
			Headers: model.JSONHeader(tsk.Headers),
			Body:    tsk.Body,
		},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task_id": id.String()})
}

func GetTaskRestServer(cfg config.ApiConfig, taskRepo db.TaskRepository) taskAPIHandler {
	router := gin.Default()

	result := taskAPIHandler{
		addr:    fmt.Sprintf(":%s", cfg.TaskAPIPort),
		handler: router,
		repo:    taskRepo,
	}

	router.GET(taskByIDEndpoint, result.getTaskByIDHandler)
	router.POST(taskEndpoint, result.addTaskHandler)

	return result

}
