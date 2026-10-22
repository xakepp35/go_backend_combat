package memory

import (
	"context"
	"jobflow/pkg"
	"sort"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	tasks map[string]pkg.Task
}

var _ pkg.TaskRepository = (*Store)(nil)

func New() *Store {
	return &Store{tasks: make(map[string]pkg.Task)}
}

func (s *Store) Insert(ctx context.Context, req pkg.InsertTaskRequest) (pkg.InsertTaskResponse, error) {
	if err := req.Task.Validate(); err != nil {
		return pkg.InsertTaskResponse{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return pkg.InsertTaskResponse{}, err
	}
	if _, exists := s.tasks[req.Task.ID]; exists {
		return pkg.InsertTaskResponse{}, pkg.ErrConflict
	}
	s.tasks[req.Task.ID] = req.Task
	return pkg.InsertTaskResponse{Task: req.Task}, nil
}

func (s *Store) List(ctx context.Context, _ pkg.ListTasksRequest) (pkg.ListTasksResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return pkg.ListTasksResponse{}, err
	}
	tasks := make([]pkg.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return pkg.ListTasksResponse{Tasks: tasks}, nil
}
