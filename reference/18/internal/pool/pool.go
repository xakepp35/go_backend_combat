package pool

import (
	"context"
	"errors"
	"sync"
)

// The first pool drains a finite input; cancellation is added on week 19.
func Run(ctx context.Context, workers int, input <-chan string, handle func(context.Context, string)) error {
	if workers < 1 {
		return errors.New("workers must be positive")
	}
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for value := range input {
				handle(ctx, value)
			}
		}()
	}
	wg.Wait()
	return nil
}
