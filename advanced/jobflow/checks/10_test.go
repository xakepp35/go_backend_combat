package coursecheck_test

import (
	"context"
	"jobflow/internal/pool"
	"sync/atomic"
	"testing"
	"time"
)

func Test10BoundedWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	jobs := make(chan func(context.Context), 8)
	gate := make(chan struct{})
	started := make(chan struct{}, 8)
	var active, finished, tooMany atomic.Int32
	for range 8 {
		jobs <- func(ctx context.Context) {
			if active.Add(1) > 2 {
				tooMany.Store(1)
			}
			started <- struct{}{}
			select {
			case <-ctx.Done():
			case <-gate:
			}
			active.Add(-1)
			finished.Add(1)
		}
	}
	close(jobs)
	done := make(chan error, 1)
	go func() { done <- pool.Run(ctx, 2, jobs) }()
	for range 2 {
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatal("workers did not start")
		}
	}
	close(gate)
	if err := <-done; err != nil || finished.Load() != 8 || tooMany.Load() != 0 {
		t.Fatal("pool bounds/drain", err, finished.Load())
	}
}
