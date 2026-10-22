package main

import (
	"context"
	"jobflow/internal/api"
	"jobflow/internal/postgres"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	s, err := postgres.Open(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	if err = s.Migrate(context.Background(), "migrations"); err != nil {
		log.Fatal(err)
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	server := &http.Server{Addr: addr, Handler: api.New(s), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
