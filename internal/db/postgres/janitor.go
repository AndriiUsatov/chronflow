package postgres

import (
	"context"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/AndriiUsatov/chronflow/internal/janitor"
	"github.com/AndriiUsatov/chronflow/internal/model"
)

type postgresJanitor struct {
	repo db.TaskRepository
}

func NewPostgresJanitor(repo db.TaskRepository) postgresJanitor {
	return postgresJanitor{
		repo: repo,
	}
}

func (janitor postgresJanitor) Run(ctx context.Context, metrics *janitor.JanitorMetrics) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			schduledReovered, _ := janitor.repo.RecoverTasks(ctx, model.Scheduled)
			failedRecovered, _ := janitor.repo.RecoverTasks(ctx, model.Failed)
			metrics.TaskRecovered.Add(float64(schduledReovered + failedRecovered))
			timer.Reset(10 * time.Minute)
		}
	}
}
