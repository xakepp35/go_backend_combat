package main

import (
	"jobflow/internal/client"
	"jobflow/internal/postgres"
	"jobflow/internal/process"
	"log/slog"
	"time"
)

func main() {
	ctx, cancel := process.Context()
	defer cancel()
	c := client.New(process.Env("PUB_URL", "http://localhost:8080"))
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for ctx.Err() == nil {
		var out postgres.GCReport
		if err := c.Do(ctx, "POST", "/internal/gc", nil, &out); err != nil {
			slog.Warn("cleanup deferred", "error", err)
		} else {
			slog.Info("cleanup complete", "removed", out.Removed)
			for _, t := range out.Tasks {
				slog.Info("task reclaimed", "workflow", t.Workflow, "task", t.ID)
			}
		}
		select {
		case <-ctx.Done():
		case <-ticker.C:
		}
	}
}
