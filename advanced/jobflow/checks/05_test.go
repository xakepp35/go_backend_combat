package coursecheck_test

import (
	"context"
	"jobflow/internal/app"
	"jobflow/internal/memory"
	"jobflow/internal/task"
	"testing"
)

type spyRepo struct {
	*memory.Repository
	writes int
}

func (r *spyRepo) Create(ctx context.Context, v task.Task) error {
	r.writes++
	return r.Repository.Create(ctx, v)
}
func Test05ConsumerBoundary(t *testing.T) {
	r := &spyRepo{Repository: memory.New()}
	s := app.New(r)
	if err := s.Create(context.Background(), task.New("", "echo", nil)); err == nil || r.writes != 0 {
		t.Fatal("invalid crossed repository boundary")
	}
	if err := s.Create(context.Background(), task.New("a", "echo", nil)); err != nil || r.writes != 1 {
		t.Fatal("service ignored injected repository", err)
	}
}
