package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/AndriiUsatov/chronflow/internal/queue"
)

// Listens RDB events and Check `RunAt` property to know when to pull tasks from RDB
// Pulls tasks from RDB and pushed it to NATS

type TaskScheduler struct {
	Repo     db.TaskRepository
	Notifier TaskNotifier
	Queue    queue.TaskQueue
}

func (scheduler TaskScheduler) Run(ctx context.Context) error {
	taskNotificationChan, err := scheduler.Notifier.Listen(ctx)
	if err != nil {
		return err
	}
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-taskNotificationChan:
			scheduler.processAndReschedule(ctx, timer)
		case <-timer.C:
			scheduler.processAndReschedule(ctx, timer)
		}
	}
}

func (scheduler TaskScheduler) processAndReschedule(ctx context.Context, timer *time.Timer) error {
	if err := scheduler.processTasks(ctx); err != nil {
		return err
	}

	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	timeNext, err := scheduler.Repo.FindNextRunAtTask(ctx)
	if err != nil {
		timer.Reset(time.Minute)
		return nil
	}

	next := time.Until(timeNext)
	if next < (time.Millisecond * 10) {
		next = 10 * time.Millisecond
	}
	timer.Reset(next)

	return nil
}

func (scheduler TaskScheduler) processTasks(ctx context.Context) error {
	// TODO: Consider limit configuration
	tasks, err := scheduler.Repo.FetchReadyTasks(ctx, 10)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := scheduler.Queue.PublishTask(ctx, task.ToProto()); err != nil {
			log.Printf("Error on publishing message: %s. Error: %v", task.ID, err)
			return err
		}
	}

	return nil
}
