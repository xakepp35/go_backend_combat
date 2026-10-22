// Package app defines dependencies where their consumers live.
package app

import (
	"context"
	"jobflow/internal/task"
)

type Repository interface {
	Create(context.Context, task.Task) error
	Get(context.Context, string) (task.Task, error)
	List(context.Context) ([]task.Task, error)
	Update(context.Context, task.Task) error
	Delete(context.Context, string) error
}
type Service struct{ repo Repository }

func New(r Repository) *Service { return &Service{repo: r} }
func (s *Service) Create(ctx context.Context, t task.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, t)
}
func (s *Service) Get(ctx context.Context, id string) (task.Task, error) { return s.repo.Get(ctx, id) }
func (s *Service) List(ctx context.Context) ([]task.Task, error)         { return s.repo.List(ctx) }
func (s *Service) Update(ctx context.Context, t task.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, t)
}
func (s *Service) Delete(ctx context.Context, id string) error { return s.repo.Delete(ctx, id) }
