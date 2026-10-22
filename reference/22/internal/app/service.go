package app

import (
	"context"
	"jobflow/pkg"
)

type Service struct{ repository pkg.TaskRepository }

var _ pkg.TaskService = (*Service)(nil)

func New(repository pkg.TaskRepository) *Service { return &Service{repository: repository} }

func (s *Service) Create(ctx context.Context, req pkg.CreateTaskRequest) (pkg.CreateTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	task := pkg.Task{ID: req.ID, Title: req.Title, Priority: req.Priority, Done: req.Done}
	if err := task.Validate(); err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	saved, err := s.repository.Insert(ctx, pkg.InsertTaskRequest{Task: task})
	if err != nil {
		return pkg.CreateTaskResponse{}, err
	}
	return pkg.CreateTaskResponse{Task: saved.Task}, nil
}
func (s *Service) Get(ctx context.Context, req pkg.GetTaskRequest) (pkg.GetTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.GetTaskResponse{}, err
	}
	return s.repository.Get(ctx, req)
}
func (s *Service) List(ctx context.Context, req pkg.ListTasksRequest) (pkg.ListTasksResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.ListTasksResponse{}, err
	}
	return s.repository.List(ctx, req)
}
func (s *Service) Update(ctx context.Context, req pkg.UpdateTaskRequest) (pkg.UpdateTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.UpdateTaskResponse{}, err
	}
	if err := req.Task.Validate(); err != nil {
		return pkg.UpdateTaskResponse{}, err
	}
	return s.repository.Update(ctx, req)
}
func (s *Service) Delete(ctx context.Context, req pkg.DeleteTaskRequest) (pkg.DeleteTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.DeleteTaskResponse{}, err
	}
	return s.repository.Delete(ctx, req)
}
func (s *Service) Complete(ctx context.Context, req pkg.CompleteTaskRequest) (pkg.CompleteTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		return pkg.CompleteTaskResponse{}, err
	}
	return s.repository.Complete(ctx, req)
}
