package main

import (
	"jobflow/internal/api"
	"jobflow/internal/postgres"
	"jobflow/internal/process"
	"log/slog"
	"os"
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
	m := api.New(s)
	api.RegisterEngine(m, s)

	err = api.RunHTTP(ctx, process.Env("HTTP_ADDR", ":8080"), m)
	cancel()
	if err != nil {
		slog.Error("http server", "error", err)
		os.Exit(1)
	}
}
