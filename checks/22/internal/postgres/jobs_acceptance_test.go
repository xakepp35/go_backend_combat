package postgres

import (
	"context"
	"errors"
	"jobflow/internal/jobs"
	"jobflow/pkg"
	"testing"
	"time"
)

func TestJobsAndRetry(t *testing.T) {
	ctx := context.Background()
	s, _, _ := isolatedStore(t)
	first, created, err := s.CreateJob(ctx, jobs.Job{ID: "j", Text: "чай хлеб молоко"})
	if err != nil || !created || first.Status != "queued" {
		t.Fatal(first, created, err)
	}
	if _, created, err = s.CreateJob(ctx, jobs.Job{ID: "j", Text: first.Text}); err != nil || created {
		t.Fatal(created, err)
	}
	if _, _, err = s.CreateJob(ctx, jobs.Job{ID: "j", Text: "other"}); !errors.Is(err, pkg.ErrConflict) {
		t.Fatal(err)
	}
	for attempt := 1; attempt <= 3; attempt++ {
		var j jobs.Job
		var found bool
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			j, found, err = s.Claim(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if found {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !found || j.Attempts != attempt {
			t.Fatal(j, found)
		}
		if err = s.Finish(ctx, j, "", jobs.ErrTemporary); err != nil {
			t.Fatal(err)
		}
	}
	final, err := s.GetJob(ctx, "j")
	if err != nil || final.Status != "failed" || final.Attempts != 3 {
		t.Fatal(final, err)
	}
}
func TestWorkersAndRecovery(t *testing.T) {
	ctx := context.Background()
	s, _, _ := isolatedStore(t)
	_, _, err := s.CreateJob(ctx, jobs.Job{ID: "j", Text: "чай хлеб молоко"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok, err := s.Claim(ctx); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if err = s.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	run, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- jobs.Run(run, s, 2, jobs.CountWords) }()
	deadline := time.Now().Add(3 * time.Second)
	var j jobs.Job
	for time.Now().Before(deadline) {
		j, err = s.GetJob(ctx, "j")
		if err != nil {
			t.Fatal(err)
		}
		if j.Status == "done" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("workers не остановились")
	}
	if j.Status != "done" || j.Result != "3" {
		t.Fatal(j)
	}
}
