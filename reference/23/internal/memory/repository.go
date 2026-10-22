package memory

import (
	"context"
	"jobflow/pkg"
)

var _ pkg.TaskRepository = (*Store)(nil)

// Port methods translate object contracts to private storage operations.
func (s *Store) Insert(ctx context.Context, req pkg.InsertTaskRequest) (pkg.InsertTaskResponse, error) {
	if err := s.create(ctx, req.Task); err != nil {
		return pkg.InsertTaskResponse{}, err
	}
	return pkg.InsertTaskResponse{Task: req.Task}, nil
}
func (s *Store) Get(ctx context.Context, req pkg.GetTaskRequest) (pkg.GetTaskResponse, error) {
	task, err := s.get(ctx, req.ID)
	if err != nil {
		return pkg.GetTaskResponse{}, err
	}
	return pkg.GetTaskResponse{Task: task}, nil
}
func (s *Store) List(ctx context.Context, req pkg.ListTasksRequest) (pkg.ListTasksResponse, error) {
	tasks, err := s.list(ctx)
	if err != nil {
		return pkg.ListTasksResponse{}, err
	}
	return pkg.ListTasksResponse{Tasks: tasks}, nil
}
func (s *Store) Update(ctx context.Context, req pkg.UpdateTaskRequest) (pkg.UpdateTaskResponse, error) {
	if err := s.update(ctx, req.Task); err != nil {
		return pkg.UpdateTaskResponse{}, err
	}
	return pkg.UpdateTaskResponse{Task: req.Task}, nil
}
func (s *Store) Delete(ctx context.Context, req pkg.DeleteTaskRequest) (pkg.DeleteTaskResponse, error) {
	if err := s.delete(ctx, req.ID); err != nil {
		return pkg.DeleteTaskResponse{}, err
	}
	return pkg.DeleteTaskResponse{}, nil
}
func (s *Store) Complete(ctx context.Context, req pkg.CompleteTaskRequest) (pkg.CompleteTaskResponse, error) {
	task, err := s.complete(ctx, req.ID)
	if err != nil {
		return pkg.CompleteTaskResponse{}, err
	}
	return pkg.CompleteTaskResponse{Task: task}, nil
}
