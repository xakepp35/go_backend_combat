package main

import (
	"jobflow/internal/api"
	"jobflow/internal/memory"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	s := &http.Server{Addr: addr, Handler: api.New(memory.New()), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(s.ListenAndServe())
}
