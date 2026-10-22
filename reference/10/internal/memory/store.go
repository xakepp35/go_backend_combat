package memory

import (
	"jobflow/pkg"
	"sort"
)

type Store struct{ values map[string]pkg.Task }

func New() *Store { return &Store{values: make(map[string]pkg.Task)} }
func (s *Store) Create(t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := s.values[t.ID]; ok {
		return pkg.ErrConflict
	}
	s.values[t.ID] = t
	return nil
}
func (s *Store) Get(id string) (pkg.Task, error) {
	t, ok := s.values[id]
	if !ok {
		return pkg.Task{}, pkg.ErrNotFound
	}
	return t, nil
}
func (s *Store) List() ([]pkg.Task, error) {
	out := make([]pkg.Task, 0, len(s.values))
	for _, t := range s.values {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Store) Update(t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := s.values[t.ID]; !ok {
		return pkg.ErrNotFound
	}
	s.values[t.ID] = t
	return nil
}
func (s *Store) Delete(id string) error {
	if _, ok := s.values[id]; !ok {
		return pkg.ErrNotFound
	}
	delete(s.values, id)
	return nil
}
func (s *Store) Complete(id string) (pkg.Task, error) {
	t, ok := s.values[id]
	if !ok {
		return pkg.Task{}, pkg.ErrNotFound
	}
	t.Complete()
	s.values[id] = t
	return t, nil
}
