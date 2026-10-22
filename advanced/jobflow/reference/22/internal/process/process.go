package process

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Context() (context.Context, context.CancelFunc) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
func Env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: Env("REDIS_ADDR", "localhost:16379"), DialTimeout: 1e9, ReadTimeout: 2e9, WriteTimeout: 2e9, MaxRetries: 0})
}
func Workers() int {
	n, err := strconv.Atoi(Env("WORKER_COUNT", "2"))
	if err != nil || n < 1 || n > 32 {
		slog.Error("WORKER_COUNT must be 1..32")
		os.Exit(1)
	}
	return n
}
