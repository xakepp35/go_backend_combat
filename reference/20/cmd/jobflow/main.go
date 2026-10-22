package main

import (
	"context"
	"errors"
	"fmt"
	"jobflow/internal/api"
	"jobflow/internal/app"
	"jobflow/internal/jobs"
	"jobflow/internal/postgres"
	"jobflow/pkg"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	startup, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	store, err := postgres.Open(startup, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer store.Close()
	if err := store.Migrate(startup); err != nil {
		return err
	}
	mux := api.New(app.New(store))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		probe, done := context.WithTimeout(r.Context(), time.Second)
		defer done()
		if err := store.Ping(probe); err != nil {
			api.WriteJSON(w, 503, pkg.HealthResponse{Status: "unavailable"})
			return
		}
		api.WriteJSON(w, 200, pkg.HealthResponse{Status: "ok"})
	})
	var workers sync.WaitGroup
	if err := store.Recover(startup); err != nil {
		return err
	}
	cancel()
	jobs.Register(mux, store)
	workers.Add(1)
	go func() {
		defer workers.Done()
		if err := jobs.Run(ctx, store, 2, jobs.CountWords); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("workers stopped", "error", err)
			stop()
		}
	}()
	server := &http.Server{Addr: httpAddress(), Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	failures := make(chan error, 1)
	go func() { failures <- server.ListenAndServe() }()
	slog.Info("server started", "address", server.Addr)
	var serveErr error
	select {
	case <-ctx.Done():
	case serveErr = <-failures:
	}
	cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	shutdownErr := server.Shutdown(cleanup)
	if shutdownErr != nil {
		_ = server.Close()
	}
	stop()
	workers.Wait()
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return serveErr
	}
	return shutdownErr
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func httpAddress() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8080"
}
