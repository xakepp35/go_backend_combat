package main

import (
	"context"
	"jobflow/internal/api"
	"jobflow/internal/app"
	"jobflow/internal/postgres"
	"jobflow/pkg"
	"log"
	"net/http"
	"os"
	"time"
)

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	store, err := postgres.Open(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer store.Close()
	if err := store.Migrate(ctx); err != nil {
		return err
	}
	cancel()
	mux := api.New(app.New(store))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		probe, stop := context.WithTimeout(r.Context(), time.Second)
		defer stop()
		if err := store.Ping(probe); err != nil {
			api.WriteJSON(w, 503, pkg.HealthResponse{Status: "unavailable"})
			return
		}
		api.WriteJSON(w, 200, pkg.HealthResponse{Status: "ok"})
	})
	s := &http.Server{Addr: httpAddress(), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Println("http://localhost:8080")
	return s.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func httpAddress() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8080"
}
