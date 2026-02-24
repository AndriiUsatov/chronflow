package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type pgTaskRepository struct {
	Db         *sqlx.DB
	taskSchema string
	taskTable  string
}

func NewPGRepositoryWithRetry(url string, port int, username, password, database, sslMode, taskSchema, taskTable string, ctx context.Context, retryFor time.Duration) (pgTaskRepository, error) {
	globalTimer := time.NewTimer(retryFor)
	sleep := time.Second
	sleepMax := time.Minute
	timer := time.NewTimer(0)

	defer globalTimer.Stop()
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return pgTaskRepository{}, ctx.Err()
		case <-globalTimer.C:
			return NewPGRepository(url, port, username, password, database, sslMode, taskSchema, taskTable)
		case <-timer.C:
			repo, err := NewPGRepository(url, port, username, password, database, sslMode, taskSchema, taskTable)
			if err == nil {
				return repo, nil
			}
			sleep *= 2
			if sleep > sleepMax {
				sleep = sleepMax
			}
			timer.Reset(sleep)
		}
	}

}

func NewPGRepository(url string, port int, username, password, database, sslMode, taskSchema, taskTable string) (pgTaskRepository, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", url, port, username, password, database, sslMode),
	)

	if err != nil {
		return pgTaskRepository{}, err
	}

	err = db.Ping()
	if err != nil {
		return pgTaskRepository{}, err
	}

	return pgTaskRepository{Db: db, taskSchema: taskSchema, taskTable: taskTable}, nil
}
