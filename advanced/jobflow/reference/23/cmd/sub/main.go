package main

import (
	"context"
	"jobflow/internal/bus"
	"jobflow/internal/client"
	"jobflow/internal/postgres"
	"jobflow/internal/process"
	"jobflow/internal/projection"
	"log/slog"
	"strings"
	"sync"
	"time"
)

func rebuild(ctx context.Context, p projection.Projection, c *client.Client) error {
	after := int64(0)
	for {
		events, err := c.Events(ctx, after)
		if err != nil {
			return err
		}
		for _, e := range events {
			if err = p.Apply(ctx, e); err != nil {
				return err
			}
			after = e.Seq
		}
		if len(events) < 500 {
			break
		}
	}
	ts, err := c.Jobs(ctx)
	if err != nil {
		return err
	}
	for _, t := range ts {
		if err = p.Apply(ctx, postgres.Event{Task: t}); err != nil {
			return err
		}
	}
	slog.Info("projection rebuilt", "tasks", len(ts))
	return nil
}
func main() {
	ctx, cancel := process.Context()
	defer cancel()
	r := process.Redis()
	defer r.Close()
	p := projection.Projection{Redis: r, Namespace: process.Env("JF_NAMESPACE", "jobflow")}
	c := client.New(process.Env("PUB_URL", "http://localhost:8080"))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bus.Consume(ctx, p, strings.Split(process.Env("KAFKA_BROKERS", "localhost:19092"), ","), process.Env("KAFKA_TOPIC", "job.events"), process.Env("KAFKA_GROUP", "jobflow-sub"))
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	initialized := false
	for ctx.Err() == nil {
		ids, err := r.SMembers(ctx, p.Key("tasks")).Result()
		if err == nil && (!initialized || len(ids) == 0) {
			err = rebuild(ctx, p, c)
			initialized = err == nil
			ids, _ = r.SMembers(ctx, p.Key("tasks")).Result()
		}
		if err == nil {
			for _, id := range ids {
				t, e := p.Get(ctx, id)
				if e != nil {
					initialized = false
					err = e
					break
				}
				if e = p.Ready(ctx, t, time.Now().UTC()); e != nil {
					err = e
					break
				}
			}
		}
		if err != nil && ctx.Err() == nil {
			slog.Warn("projection reconciliation deferred", "error", err)
		}
		select {
		case <-ctx.Done():
		case <-ticker.C:
		}
	}
	wg.Wait()
}
