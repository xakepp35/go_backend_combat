package memory

import (
	"context"
	"jobflow/internal/task"
	"sort"
	"sync"
)

type Repository struct {
	mu    sync.Mutex
	tasks map[string]task.Task
}

func New() *Repository { return &Repository{tasks: map[string]task.Task{}} }
func (r *Repository) Create(ctx context.Context, t task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := t.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[t.ID]; ok {
		return task.ErrConflict
	}
	t.State = task.Draft
	t.Version = 1
	r.tasks[t.ID] = task.Clone(t)
	return nil
}
func (r *Repository) Get(ctx context.Context, id string) (task.Task, error) {
	if err := ctx.Err(); err != nil {
		return task.Task{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tasks[id]
	if !ok {
		return task.Task{}, task.ErrNotFound
	}
	return task.Clone(t), nil
}
func (r *Repository) List(ctx context.Context) ([]task.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []task.Task{}
	for _, t := range r.tasks {
		out = append(out, task.Clone(t))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (r *Repository) Update(ctx context.Context, t task.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := t.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	old, ok := r.tasks[t.ID]
	if !ok {
		return task.ErrNotFound
	}
	if old.State != task.Draft {
		return task.ErrConflict
	}
	t.State = task.Draft
	t.Version = old.Version + 1
	r.tasks[t.ID] = task.Clone(t)
	return nil
}
func (r *Repository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tasks[id]
	if !ok {
		return task.ErrNotFound
	}
	if t.State != task.Draft {
		return task.ErrConflict
	}
	delete(r.tasks, id)
	return nil
}
