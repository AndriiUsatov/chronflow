package api

import (
	"context"
	"net"

	"github.com/AndriiUsatov/chronflow/internal/db"
	"github.com/AndriiUsatov/chronflow/internal/model"
	"github.com/AndriiUsatov/chronflow/internal/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type taskStatusServerImpl struct {
	pb.UnimplementedTaskServiceServer
	taskRepository db.TaskRepository
}

type taskStatusServer struct {
	listener net.Listener
	server   *grpc.Server
}

func (server taskStatusServer) ListenAndServe(isPanic bool) error {
	err := server.server.Serve(server.listener)
	if isPanic && err != nil {
		panic(err)
	}
	return err
}

func (server taskStatusServer) Close() {
	server.server.GracefulStop()
}

func NewTaskStatusServer(network, address string, repo db.TaskRepository) (taskStatusServer, error) {
	lis, err := net.Listen(network, address)
	if err != nil {
		return taskStatusServer{}, err
	}

	ser := grpc.NewServer()
	pb.RegisterTaskServiceServer(ser, &taskStatusServerImpl{
		taskRepository: repo,
	})

	return taskStatusServer{
		listener: lis,
		server:   ser,
	}, nil
}

func (server *taskStatusServerImpl) UpdateTaskStatus(ctx context.Context, req *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateTaskStatusResponse{Success: false}, err
	}

	err = server.taskRepository.UpdateTaskStatus(ctx, id, model.TaskStatus(req.Status), req.ErrorMessage)

	return &pb.UpdateTaskStatusResponse{Success: (err == nil)}, err
}
