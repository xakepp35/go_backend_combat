package main

import (
	"jobflow/internal/api"
	"log"
	"net/http"
	"os"
	"time"
)

func run() error {
	server := &http.Server{Addr: httpAddress(), Handler: api.New(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("starting HTTP on %s", server.Addr)
	return server.ListenAndServe()
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
