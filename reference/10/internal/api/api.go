package api

import (
	"fmt"
	"jobflow/pkg"
	"net/http"
	"strings"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			fmt.Fprintln(w, "Привет из "+pkg.AppName)
			return
		}
		fmt.Fprintf(w, "Привет, %s! Это %s.\n", name, pkg.AppName)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, pkg.AppName+" — список дел. Учимся превращать просьбы в программы.")
	})
	return mux
}
