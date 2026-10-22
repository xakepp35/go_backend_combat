package app

import (
	"context"
	"jobflow/pkg"
)

type Service struct {
	repository pkg.TaskRepository
}

var _ pkg.TaskService = (*Service)(nil)

func New(repository pkg.TaskRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, req pkg.CreateTaskRequest) (pkg.CreateTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	task := pkg.Task{ID: req.ID, Title: req.Title, Priority: req.Priority}
	if err := task.Validate(); err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	saved, err := s.repository.Insert(ctx, pkg.InsertTaskRequest{Task: task})
	if err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	return pkg.CreateTaskResponse{Task: saved.Task}, nil
}

func (s *Service) List(ctx context.Context, req pkg.ListTasksRequest) (pkg.ListTasksResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.ListTasksResponse{}, err
	}
	return s.repository.List(ctx, req)
}
