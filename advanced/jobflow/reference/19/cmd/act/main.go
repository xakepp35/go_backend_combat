package main

import (
	"jobflow/internal/client"
	"jobflow/internal/process"
	"jobflow/internal/projection"
	"jobflow/internal/runner"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := process.Context()
	defer cancel()
	r := process.Redis()
	defer r.Close()
	p := projection.Projection{Redis: r, Namespace: process.Env("JF_NAMESPACE", "jobflow")}
	c := client.New(process.Env("PUB_URL", "http://localhost:8080"))
	if err := runner.Run(ctx, p, c, process.Workers()); err != nil && ctx.Err() == nil {
		slog.Error("runner", "error", err)
		os.Exit(1)
	}
}
