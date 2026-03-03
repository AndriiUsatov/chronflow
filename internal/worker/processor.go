package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/AndriiUsatov/chronflow/internal/pb"
	"github.com/AndriiUsatov/chronflow/internal/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	retries              = 3
	userAgentHeaderKey   = "User-Agent"
	userAgentHeaderValue = "ChronFlow-Worker/1.0"
)

type TaskWorker struct {
	queue         queue.TaskQueue
	httpClient    *http.Client
	serviceClient pb.TaskServiceClient
	grpcConn      *grpc.ClientConn
	metrics       *WorkerMetrics
}

func (worker TaskWorker) Close() {
	worker.grpcConn.Close()
}

func NewTaskWorker(ctx context.Context, queue queue.TaskQueue, target string, retryFor time.Duration, metrics *WorkerMetrics) (TaskWorker, error) {
	globalTimer := time.NewTimer(retryFor)
	defer globalTimer.Stop()

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return TaskWorker{}, ctx.Err()
		case <-globalTimer.C:
			return tryConnect(target, queue, timer, metrics)
		case <-timer.C:
			res, err := tryConnect(target, queue, timer, metrics)
			if err == nil {
				return res, nil
			}
		}
	}
}

func tryConnect(target string, queue queue.TaskQueue, timer *time.Timer, metrics *WorkerMetrics) (TaskWorker, error) {
	// TODO: Replace insecure creds
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)

	if err != nil {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		//TODO: Consider sleep configuration
		timer.Reset(5 * time.Second)

		return TaskWorker{}, err
	}

	client := pb.NewTaskServiceClient(conn)

	return TaskWorker{
		queue: queue,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		serviceClient: client,
		grpcConn:      conn,
		metrics:       metrics,
	}, nil
}

func (worker TaskWorker) Run(ctx context.Context) error {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			found, err := worker.processTask(ctx, timer)
			if err != nil {
				log.Printf("Worker error: %v\n", err)
				timer.Reset(2 * time.Second)
			} else if found {
				timer.Reset(0)
			} else {
				timer.Reset(time.Second)
			}
		}
	}

}

func (worker TaskWorker) processTask(ctx context.Context, timer *time.Timer) (bool, error) {

	task, err := worker.queue.ConsumeTask(ctx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return false, nil
		}
		return false, fmt.Errorf("Queue consume error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		task.Method,
		task.Url,
		bytes.NewBuffer(task.Body),
	)

	if err != nil {
		worker.handleTaskResult(ctx, task.Id, model.Failed, err.Error(), retries)
		return true, nil
	}

	for k, v := range task.Headers {
		for _, val := range v.Values {
			req.Header.Add(k, val)
		}
	}
	req.Header.Add(userAgentHeaderKey, userAgentHeaderValue)

	res, err := worker.httpClient.Do(req)

	if err != nil {
		worker.handleTaskResult(ctx, task.Id, model.Failed, err.Error(), retries)
		return true, nil
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		worker.handleTaskResult(ctx, task.Id, model.Failed, fmt.Sprintf("Error on request, status code: %d", res.StatusCode), retries)
		return true, nil
	}

	worker.handleTaskResult(ctx, task.Id, model.Completed, "", retries)
	return true, nil
}

func (worker TaskWorker) handleTaskResult(ctx context.Context, id string, status model.TaskStatus, errMsg string, retries int) {
	for i := 0; i < retries; i++ {
		res, err := worker.serviceClient.UpdateTaskStatus(ctx, &pb.UpdateTaskStatusRequest{
			Id:           id,
			Status:       pb.Status(status),
			ErrorMessage: errMsg,
		})

		if res.Success && err == nil {
			if status == model.Failed {
				worker.metrics.taskProcessed.WithLabelValues(statusFail).Inc()
			} else {
				worker.metrics.taskProcessed.WithLabelValues(statusSuccess).Inc()
			}
			return
		}
	}
	worker.metrics.taskProcessed.WithLabelValues(statusFail).Inc()
}
