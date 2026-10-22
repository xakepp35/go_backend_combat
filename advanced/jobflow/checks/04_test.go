package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/memory"
	"jobflow/internal/task"
	"testing"
)

func Test04InvalidDoesNotMutate(t *testing.T) {
	r := memory.New()
	ctx := context.Background()
	for _, v := range []task.Task{task.New("", "echo", nil), task.New("a", "", nil), {ID: "a", Kind: "echo", NBF: 10, EXP: 10}} {
		if err := r.Create(ctx, v); !errors.Is(err, task.ErrInvalid) {
			t.Fatal("invalid not recognized", err)
		}
	}
	list, _ := r.List(ctx)
	if len(list) != 0 {
		t.Fatal("invalid input changed state")
	}
	if _, err := r.Get(ctx, "missing"); !errors.Is(err, task.ErrNotFound) {
		t.Fatal(err)
	}
}
