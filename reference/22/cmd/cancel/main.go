package main

import (
	"context"
	"fmt"
	"time"
)

func wait(ctx context.Context) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	fmt.Println(wait(ctx))
}
