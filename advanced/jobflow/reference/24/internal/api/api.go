// Package api translates HTTP requests to application calls.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"jobflow/internal/app"
	"jobflow/internal/task"
	"log/slog"
	"net/http"
)

func Write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, task.ErrNotFound):
		status = 404
	case errors.Is(err, task.ErrInvalid):
		status = 400
	case errors.Is(err, task.ErrConflict), errors.Is(err, task.ErrNotReady), errors.Is(err, task.ErrStale):
		status = 409
	case errors.Is(err, task.ErrExpired):
		status = 410
	case errors.Is(err, task.ErrQuota):
		status = 422
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		status = 408
	}
	message := err.Error()
	if status == 500 {
		slog.Error("request failed", "error", err)
		message = "internal error"
	}
	Write(w, status, map[string]string{"error": message})
}
func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return task.ErrInvalid
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return task.ErrInvalid
	}
	return nil
}
func New(repo app.Repository) *http.ServeMux {
	s := app.New(repo)
	m := http.NewServeMux()
	m.HandleFunc("POST /drafts", func(w http.ResponseWriter, r *http.Request) {
		var t task.Task
		if err := Decode(w, r, &t); err != nil {
			Error(w, err)
			return
		}
		if err := s.Create(r.Context(), t); err != nil {
			Error(w, err)
			return
		}
		Write(w, 201, map[string]string{"id": t.ID})
	})
	m.HandleFunc("GET /drafts", func(w http.ResponseWriter, r *http.Request) {
		ts, err := s.List(r.Context())
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, ts)
	})
	m.HandleFunc("GET /drafts/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, t)
	})
	m.HandleFunc("PUT /drafts/{id}", func(w http.ResponseWriter, r *http.Request) {
		var t task.Task
		if err := Decode(w, r, &t); err != nil {
			Error(w, err)
			return
		}
		t.ID = r.PathValue("id")
		if err := s.Update(r.Context(), t); err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, map[string]string{"id": t.ID})
	})
	m.HandleFunc("DELETE /drafts/{id}", func(w http.ResponseWriter, r *http.Request) {
		if err := s.Delete(r.Context(), r.PathValue("id")); err != nil {
			Error(w, err)
			return
		}
		w.WriteHeader(204)
	})
	m.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { Write(w, 200, map[string]string{"status": "ok"}) })
	return m
}
