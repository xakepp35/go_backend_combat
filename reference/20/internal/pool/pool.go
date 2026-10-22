package pool

import (
	"context"
	"errors"
	"sync"
)

// Run waits for all workers; the producer retains ownership of closing input.
// Each handler must observe ctx when it blocks.
func Run(ctx context.Context, workers int, input <-chan string, handle func(context.Context, string)) error {
	if workers < 1 {
		return errors.New("workers must be positive")
	}
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				case value, ok := <-input:
					if !ok {
						return
					}
					handle(ctx, value)
				}
			}
		}()
	}
	wg.Wait()
	return ctx.Err()
}
