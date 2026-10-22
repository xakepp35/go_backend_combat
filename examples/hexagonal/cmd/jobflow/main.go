package main

import (
	"jobflow/internal/api"
	"jobflow/internal/app"
	"jobflow/internal/memory"
	"log"
	"net/http"
	"time"
)

func main() {
	store := memory.New()
	service := app.New(store)
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           api.New(service),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("JobFlow: http://127.0.0.1:8080")
	log.Fatal(server.ListenAndServe())
}
