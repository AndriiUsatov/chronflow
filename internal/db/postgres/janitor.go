package postgres

import (
	"context"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/db"
)

type postgresJanitor struct {
	repo db.TaskRepository
}

func NewPostgresJanitor(repo db.TaskRepository) postgresJanitor {
	return postgresJanitor{
		repo: repo,
	}
}

func (janitor postgresJanitor) Run(ctx context.Context) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			janitor.repo.RecoverStuckTasks(ctx)
			timer.Reset(10 * time.Minute)
		}
	}
}
