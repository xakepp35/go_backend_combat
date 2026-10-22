package main

import (
	"jobflow/internal/api"
	"jobflow/internal/app"
	"jobflow/internal/memory"
	"jobflow/pkg"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := api.New(app.New(memory.New()))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { api.WriteJSON(w, 200, pkg.HealthResponse{Status: "ok"}) })
	s := &http.Server{Addr: httpAddress(), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Println("http://localhost:8080")
	if err := s.ListenAndServe(); err != nil {
		log.Print(err)
	}
}

func httpAddress() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8080"
}
