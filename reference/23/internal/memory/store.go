package memory

import (
	"context"
	"jobflow/pkg"
	"sort"
	"sync"
)

type Store struct {
	values map[string]pkg.Task
	mu     sync.Mutex
}

func New() *Store { return &Store{values: make(map[string]pkg.Task)} }
func (s *Store) create(ctx context.Context, t pkg.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := s.values[t.ID]; ok {
		return pkg.ErrConflict
	}
	s.values[t.ID] = t
	return nil
}
func (s *Store) get(ctx context.Context, id string) (pkg.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return pkg.Task{}, err
	}
	t, ok := s.values[id]
	if !ok {
		return pkg.Task{}, pkg.ErrNotFound
	}
	return t, nil
}
func (s *Store) list(ctx context.Context) ([]pkg.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := make([]pkg.Task, 0, len(s.values))
	for _, t := range s.values {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Store) update(ctx context.Context, t pkg.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := s.values[t.ID]; !ok {
		return pkg.ErrNotFound
	}
	s.values[t.ID] = t
	return nil
}
func (s *Store) delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, ok := s.values[id]; !ok {
		return pkg.ErrNotFound
	}
	delete(s.values, id)
	return nil
}
func (s *Store) complete(ctx context.Context, id string) (pkg.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return pkg.Task{}, err
	}
	t, ok := s.values[id]
	if !ok {
		return pkg.Task{}, pkg.ErrNotFound
	}
	t.Complete()
	s.values[id] = t
	return t, nil
}
