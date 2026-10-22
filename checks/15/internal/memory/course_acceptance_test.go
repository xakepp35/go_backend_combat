package memory

import (
	"context"
	"errors"
	"jobflow/pkg"
	"testing"
)

func TestRepositoryContract(t *testing.T) {
	ctx := context.Background()
	s := New()
	a, _ := pkg.New("a", "чай", 1)
	if _, err := s.Insert(ctx, pkg.InsertTaskRequest{Task: a}); err != nil {
		t.Fatal(err)
	}
	b := a
	b.Title = "кофе"
	if _, err := s.Insert(ctx, pkg.InsertTaskRequest{Task: b}); !errors.Is(err, pkg.ErrConflict) {
		t.Fatal("дубликат")
	}
	got, err := s.Get(ctx, pkg.GetTaskRequest{ID: "a"})
	if err != nil || got.Task.Title != "чай" {
		t.Fatal("первое значение потеряно", err)
	}
	if _, err := s.Update(ctx, pkg.UpdateTaskRequest{Task: pkg.Task{ID: "z", Title: "x", Priority: 1}}); !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("неизвестное изменение")
	}
	done, err := s.Complete(ctx, pkg.CompleteTaskRequest{ID: "a"})
	if err != nil || !done.Task.Done {
		t.Fatal(err)
	}
	if _, err = s.Delete(ctx, pkg.DeleteTaskRequest{ID: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Get(ctx, pkg.GetTaskRequest{ID: "a"}); !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("удаление")
	}
}
