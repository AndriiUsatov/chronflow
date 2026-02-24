package queue

import (
	"context"

	"github.com/AndriiUsatov/chronflow/internal/pb"
)

type TaskQueue interface {
	Close()
	PublishTask(ctx context.Context, task *pb.ProtoTask) error
	ConsumeTask(ctx context.Context) (*pb.ProtoTask, error)
}
