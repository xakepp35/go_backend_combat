package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/postgres"
	"jobflow/internal/task"
	"testing"
	"time"
)

func Test22SharedCASRetention(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	jobs := submit(t, s, []byte("same"), []byte("same"))
	key := jobs[0].Key
	if key != jobs[1].Key {
		t.Fatal("CAS not shared")
	}
	for _, job := range jobs {
		a, err := s.Claim(ctx, job.ID, time.Unix(1001, 0), 5*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = s.Complete(ctx, a.ID, postgres.Result{Token: a.Token, Body: []byte("same")}, time.Unix(1002, 0)); err != nil {
			t.Fatal(err)
		}
		if _, err = s.GC(ctx, time.Unix(2000, 0)); err != nil {
			t.Fatal(err)
		}
		if _, err = s.Payload(ctx, key); err != nil {
			t.Fatal("GC deleted undelivered/shared payload", err)
		}
	}
	delivered(t, s)
	n, err := s.GC(ctx, time.Unix(2000, 0))
	if err != nil || n != 1 {
		t.Fatal("GC did not reclaim shared blob once", n, err)
	}
	if _, err = s.Payload(ctx, key); !errors.Is(err, task.ErrNotFound) {
		t.Fatal(err)
	}
	got, err := s.Job(ctx, jobs[0].ID)
	if err != nil || got.State != task.Succeeded {
		t.Fatal("GC destroyed identity/state", got, err)
	}
}

func Test22LiveDependentProtectsPayload(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	ts, err := s.Submit(ctx, dependent(t, frames(t, []byte("source"), []byte("waiting"))), time.Unix(1000, 0))
	if err != nil {
		t.Fatal(err)
	}
	a, err := s.Claim(ctx, ts[0].ID, time.Unix(1001, 0), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Complete(ctx, a.ID, postgres.Result{Token: a.Token, Body: []byte("result")}, time.Unix(1001, 0)); err != nil {
		t.Fatal(err)
	}
	delivered(t, s)
	report, err := s.Collect(ctx, time.Unix(2000, 0))
	if err != nil || len(report.Tasks) != 0 {
		t.Fatal("live dependent lost source", report, err)
	}
	if _, err = s.Payload(ctx, ts[0].Key); err != nil {
		t.Fatal(err)
	}
	a, err = s.Claim(ctx, ts[1].ID, time.Unix(1002, 0), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Complete(ctx, a.ID, postgres.Result{Token: a.Token}, time.Unix(1002, 0)); err != nil {
		t.Fatal(err)
	}
	delivered(t, s)
	report, err = s.Collect(ctx, time.Unix(2000, 0))
	if err != nil || len(report.Tasks) != 2 {
		t.Fatal("eligible tasks not reported", report, err)
	}
	report, err = s.Collect(ctx, time.Unix(2000, 0))
	if err != nil || report.Removed != 0 || len(report.Tasks) != 0 {
		t.Fatal("cleanup retry changed result", report, err)
	}
}
