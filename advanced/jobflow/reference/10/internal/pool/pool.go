package pool

import (
	"context"
	"errors"
	"sync"
)

func Run(ctx context.Context, workers int, jobs <-chan func(context.Context)) error {
	if workers < 1 || workers > 256 {
		return errors.New("workers must be 1..256")
	}
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range jobs {
				f(ctx)
			}
		}()
	}
	wg.Wait()
	return nil
}
