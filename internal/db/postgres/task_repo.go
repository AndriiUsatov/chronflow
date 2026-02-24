package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/google/uuid"
)

func (repo pgTaskRepository) CreateTask(ctx context.Context, task model.Task) error {
	_, err := repo.Db.NamedExecContext(
		ctx,
		fmt.Sprintf(`
		INSERT INTO %s.%s 
			(id, url, method, headers, body, run_at)
		VALUES (:id, :url, :method, :headers, :body, :run_at)
		`, repo.taskSchema, repo.taskTable),
		&task)

	return err
}

func (repo pgTaskRepository) GetTasksByStatus(ctx context.Context, status model.TaskStatus, limit int) (res model.Tasks, err error) {
	err = repo.Db.SelectContext(
		ctx,
		&res,
		fmt.Sprintf(
			`SELECT
				id,
				url,
				method,
				headers,
				body,
				status,
				run_at,
				created,
				updated,
				retry_count,
				COALESCE(error_message, '') as error_message 
		 	FROM %s.%s WHERE status=$1 LIMIT $2`, repo.taskSchema, repo.taskTable),
		status, limit)

	return
}

func (repo pgTaskRepository) GetTasksByStatusAndRunAt(ctx context.Context, status model.TaskStatus, runAt time.Time, limit int) (res model.Tasks, err error) {
	err = repo.Db.SelectContext(
		ctx,
		&res,
		fmt.Sprintf(
			`SELECT 
				id,
				url,
				method,
				headers,
				body,
				status,
				run_at,
				created,
				updated,
				retry_count,
				COALESCE(error_message, '') as error_message	 
			FROM %s.%s WHERE status=$1 AND run_at >= $2 LIMIT $3`, repo.taskSchema, repo.taskTable),
		repo.taskSchema, repo.taskTable, status, runAt, limit)

	return
}

func (repo pgTaskRepository) GetTaskByUUID(ctx context.Context, uid uuid.UUID) (res model.Task, err error) {
	err = repo.Db.GetContext(
		ctx,
		&res,
		fmt.Sprintf(
			`SELECT 
				id,
				url,
				method,
				headers,
				body,
				status,
				run_at,
				created,
				updated,
				retry_count,
				COALESCE(error_message, '') as error_message
			FROM %s.%s WHERE id=$1`, repo.taskSchema, repo.taskTable),
		uid)

	return
}

func (repo pgTaskRepository) FetchReadyTasks(ctx context.Context, limit int) (res model.Tasks, err error) {
	err = repo.Db.SelectContext(
		ctx,
		&res,
		fmt.Sprintf(`
			UPDATE %s.%s
			SET 
				status = $1,
				updated = current_timestamp
			WHERE id IN (
			    SELECT id
			    FROM %s.%s
			    WHERE status = $2 AND current_timestamp >= run_at 
			    ORDER BY run_at ASC
			    LIMIT $3
			    FOR UPDATE SKIP LOCKED
			)
			RETURNING 
				id,
				url,
				method,
				headers,
				body,
				status,
				run_at,
				created,
				updated,
				retry_count,
				COALESCE(error_message, '') as error_message
			;`,
			repo.taskSchema, repo.taskTable, repo.taskSchema, repo.taskTable),
		model.Scheduled, model.Pending, limit)

	return
}

func (repo pgTaskRepository) FindNextRunAtTask(ctx context.Context) (res time.Time, err error) {
	err = repo.Db.GetContext(
		ctx,
		&res,
		fmt.Sprintf("SELECT MIN(run_at) FROM %s.%s LIMIT 1", repo.taskSchema, repo.taskTable),
	)

	return
}

func (repo pgTaskRepository) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus, errMsg string) error {
	_, err := repo.Db.ExecContext(
		ctx,
		fmt.Sprintf(`
			UPDATE %s.%s
			SET 
				status = $2,
				updated = current_timestamp,
				error_message = 
				CASE
					WHEN $2 = $3
					THEN $4
					ELSE NULL
				END
			WHERE
				id=$1 AND status != $2
		`, repo.taskSchema, repo.taskTable),
		id, status, model.Failed, errMsg,
	)

	return err
}

func (repo pgTaskRepository) RecoverStuckTasks(ctx context.Context) error {
	_, err := repo.Db.ExecContext(
		ctx,
		fmt.Sprintf(`
			UPDATE %s.%s
			SET
				status = $1,
				updated = current_timestamp,
				retry_count = retry_count + 1
			WHERE
				status = $2
				AND 
				updated < (current_timestamp - INTERVAL '10 minutes')
				AND
				retry_count < 5
		`, repo.taskSchema, repo.taskTable),
		model.Pending, model.Scheduled)

	return err
}
