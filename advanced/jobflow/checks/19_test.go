package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/postgres"
	"jobflow/internal/runner"
	"jobflow/internal/task"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test19AtomicClaim(t *testing.T) {
	s, _ := database(t)
	job := submit(t, s, []byte("work"))[0]
	ctx := context.Background()
	var winners atomic.Int32
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := s.Claim(ctx, job.ID, time.Unix(1001, 0), 5*time.Second); err == nil {
				winners.Add(1)
			} else if !errors.Is(err, task.ErrNotReady) {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if winners.Load() != 1 {
		t.Fatal("multiple valid attempts", winners.Load())
	}
	if _, err := s.Complete(ctx, job.ID, postgres.Result{Token: "wrong"}, time.Unix(1002, 0)); !errors.Is(err, task.ErrStale) {
		t.Fatal("foreign token accepted", err)
	}
	body, err := runner.Execute(ctx, task.New("test", "sum", []byte("[1,2,3]")), runner.Registry(), nil)
	if err != nil || string(body) != "6" {
		t.Fatal("handler", string(body), err)
	}
}
