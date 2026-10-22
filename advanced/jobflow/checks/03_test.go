package coursecheck_test

import (
	"context"
	"jobflow/internal/memory"
	"jobflow/internal/task"
	"testing"
)

func Test03CRUD(t *testing.T) {
	ctx := context.Background()
	r := memory.New()
	if err := r.Create(ctx, task.New("a", "echo", []byte("old"))); err != nil {
		t.Fatal(err)
	}
	if err := r.Create(ctx, task.New("a", "echo", nil)); err == nil {
		t.Fatal("duplicate overwrote task")
	}
	got, err := r.Get(ctx, "a")
	if err != nil || string(got.Body) != "old" {
		t.Fatal(got, err)
	}
	got.Body = []byte("new")
	if err = r.Update(ctx, got); err != nil {
		t.Fatal(err)
	}
	list, err := r.List(ctx)
	if err != nil || len(list) != 1 || string(list[0].Body) != "new" {
		t.Fatal(list, err)
	}
	if err = r.Delete(ctx, "a"); err != nil {
		t.Fatal(err)
	}
	if _, err = r.Get(ctx, "a"); err == nil {
		t.Fatal("deleted task exists")
	}
}
