package pool

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestReferenceBoundAndCompletion(t *testing.T) {
	input := make(chan string, 10)
	for i := 0; i < 10; i++ {
		input <- "x"
	}
	close(input)
	var mu sync.Mutex
	active, max, finished := 0, 0, 0
	started := make(chan struct{}, 10)
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	defer unblock()
	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), 2, input, func(ctx context.Context, s string) {
			mu.Lock()
			active++
			if active > max {
				max = active
			}
			mu.Unlock()
			started <- struct{}{}
			<-release
			mu.Lock()
			active--
			finished++
			mu.Unlock()
		})
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("не запущены два исполнителя")
		}
	}
	select {
	case <-started:
		t.Fatal("лимит превышен: третий исполнитель вошёл до освобождения первых двух")
	case <-time.After(50 * time.Millisecond):
	}
	unblock()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("pool не завершился")
	}
	mu.Lock()
	defer mu.Unlock()
	if max > 2 || active != 0 || finished != 10 {
		t.Fatal(max, active, finished)
	}
}
func TestReferenceRejectsZeroWorkers(t *testing.T) {
	if err := Run(context.Background(), 0, make(chan string), func(context.Context, string) {}); err == nil {
		t.Fatal("нужен положительный лимит")
	}
}
