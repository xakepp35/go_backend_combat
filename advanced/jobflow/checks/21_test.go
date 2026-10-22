package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/postgres"
	"jobflow/internal/runner"
	"jobflow/internal/task"
	"testing"
	"time"
)

func Test21CrashAfterEffect(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	job := submit(t, s, []byte("effect"))[0]
	a, err := s.Claim(ctx, job.ID, time.Unix(1001, 0), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	first, err := s.Effect(ctx, a.ID, a.Token, "counter", 5, time.Unix(1001, 0))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Maintain(ctx, time.Unix(1003, 0)); err != nil {
		t.Fatal(err)
	}
	b, err := s.Claim(ctx, job.ID, time.Unix(1004, 0), 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	again, err := s.Effect(ctx, b.ID, b.Token, "counter", 5, time.Unix(1004, 0))
	if err != nil || first != again {
		t.Fatal("effect duplicated", first, again, err)
	}
	if _, err = s.Effect(ctx, a.ID, a.Token, "other", 1, time.Unix(1004, 0)); !errors.Is(err, task.ErrStale) {
		t.Fatal("old attempt unfenced", err)
	}
	result := postgres.Result{Token: b.Token, Body: []byte("ok")}
	for range 2 {
		if _, err = s.Complete(ctx, b.ID, result, time.Unix(1005, 0)); err != nil {
			t.Fatal("result retry", err)
		}
	}
	result.Body = []byte("different")
	if _, err = s.Complete(ctx, b.ID, result, time.Unix(1005, 0)); !errors.Is(err, task.ErrConflict) {
		t.Fatal("result changed", err)
	}
}
func Test21PanicAndQuota(t *testing.T) {
	registry := runner.Registry()
	registry["panic"] = func(context.Context, task.Task, *runner.Budget, runner.Effects) ([]byte, error) { panic("injected") }
	if _, err := runner.Execute(context.Background(), task.New("p", "panic", nil), registry, nil); err == nil {
		t.Fatal("panic vanished")
	}
	v := task.New("s", "sum", []byte("[1,2]"))
	v.Quotas.FX = 1
	if _, err := runner.Execute(context.Background(), v, registry, nil); !errors.Is(err, task.ErrQuota) {
		t.Fatal("quota ignored", err)
	}
}

func Test21AbandonedAttemptsAreBounded(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	job := submit(t, s, []byte("abandoned"))[0]
	now := time.Unix(1001, 0)
	for attempt := 1; attempt <= 3; attempt++ {
		a, err := s.Claim(ctx, job.ID, now, time.Second)
		if err != nil || a.Attempts != attempt {
			t.Fatal("claim", attempt, a, err)
		}
		now = now.Add(2 * time.Second)
		if err = s.Maintain(ctx, now); err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Second)
	}
	got, err := s.Job(ctx, job.ID)
	if err != nil || got.State != task.Failed {
		t.Fatal("abandoned job retries forever", got, err)
	}
	if _, err = s.Claim(ctx, job.ID, now, time.Second); err == nil {
		t.Fatal("fourth attempt allowed")
	}
}
