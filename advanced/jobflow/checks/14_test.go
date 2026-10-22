package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/task"
	"testing"
	"time"
)

func Test14DAGAndSubmission(t *testing.T) {
	if err := task.ValidateDAG([]task.Task{{ID: "a", Kind: "echo", Refs: []string{"b"}}, {ID: "b", Kind: "echo", Refs: []string{"a"}}}); !errors.Is(err, task.ErrInvalid) {
		t.Fatal("cycle accepted", err)
	}
	s, _ := database(t)
	raw := dependent(t, frames(t, []byte("a"), []byte("b")))
	ts, err := s.Submit(context.Background(), raw, time.Unix(1000, 0))
	if err != nil || len(ts) != 2 || len(ts[1].Refs) != 1 || ts[1].Refs[0] != ts[0].ID {
		t.Fatal("refs not resolved", ts, err)
	}
	if _, err = s.Submit(context.Background(), raw, time.Unix(1001, 0)); err != nil {
		t.Fatal("idempotent workflow", err)
	}
}

func Test14EmptyPayloadIsNotSQLNull(t *testing.T) {
	s, _ := database(t)
	jobs := submit(t, s, nil)
	got, err := s.Job(context.Background(), jobs[0].ID)
	if err != nil || len(got.Body) != 0 {
		t.Fatal("empty CAS payload", got, err)
	}
}
