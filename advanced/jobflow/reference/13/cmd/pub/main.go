package main

import (
	"context"
	"jobflow/internal/api"
	"jobflow/internal/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s, err := postgres.Open(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	if err = s.Migrate(ctx, "migrations"); err != nil {
		log.Fatal(err)
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err = api.RunHTTP(ctx, addr, api.New(s)); err != nil {
		log.Fatal(err)
	}
}
