package db

import (
	"context"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task model.Task) error
	GetTasksByStatus(ctx context.Context, status model.TaskStatus, limit int) (model.Tasks, error)
	GetTasksByStatusAndRunAt(ctx context.Context, status model.TaskStatus, runAt time.Time, limit int) (model.Tasks, error)
	GetTaskByUUID(ctx context.Context, uid uuid.UUID) (model.Task, error)
	FetchReadyTasks(ctx context.Context, limit int) (model.Tasks, error)
	FindNextRunAtTask(ctx context.Context) (time.Time, error)
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus, errMsg string) error
	RecoverTasks(ctx context.Context, targetStatus model.TaskStatus) (int64, error)
	Ping(ctx context.Context) error
}
