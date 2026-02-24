package scheduler

import "context"

type TaskNotifier interface {
	Listen(ctx context.Context) (<-chan struct{}, error)
}
