package main

import (
	"jobflow/internal/api"
	"jobflow/internal/bus"
	"jobflow/internal/cache"
	"jobflow/internal/postgres"
	"jobflow/internal/process"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	ctx, cancel := process.Context()
	defer cancel()
	s, err := postgres.Open(ctx, process.Env("DATABASE_URL", "postgres://jobflow:jobflow@localhost:15432/jobflow?sslmode=disable"))
	if err != nil {
		slog.Error("open postgres", "error", err)
		os.Exit(1)
	}
	defer s.Close()
	if err = s.Migrate(ctx, "migrations"); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}
	r := process.Redis()
	defer r.Close()
	cas := cache.CAS{Redis: r, Source: s, Namespace: process.Env("JF_NAMESPACE", "jobflow")}
	m := api.New(s)
	api.RegisterEngine(m, s, cas, time.Hour)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		bus.Relay(ctx, s, strings.Split(process.Env("KAFKA_BROKERS", "localhost:19092"), ","), process.Env("KAFKA_TOPIC", "job.events"))
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				if err := s.Maintain(ctx, now.UTC()); err != nil && ctx.Err() == nil {
					slog.Warn("maintenance deferred", "error", err)
				}
			}
		}
	}()
	err = api.RunHTTP(ctx, process.Env("HTTP_ADDR", ":8080"), m)
	cancel()
	wg.Wait()
	if err != nil {
		slog.Error("http server", "error", err)
		os.Exit(1)
	}
}
