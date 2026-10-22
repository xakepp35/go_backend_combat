package api

import (
	"encoding/json"
	"errors"
	"io"
	"jobflow/pkg"
	"log/slog"
	"net/http"
)

// New depends on the public service port, not a concrete app.Service.
func New(service pkg.TaskService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var req pkg.CreateTaskRequest
		if err := decoder.Decode(&req); err != nil {
			writeError(w, pkg.ErrInvalid)
			return
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeError(w, pkg.ErrInvalid)
			return
		}
		res, err := service.Create(r.Context(), req)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, res)
	})
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.List(r.Context(), pkg.ListTasksRequest{})
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
	return mux
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal error"
	switch {
	case errors.Is(err, pkg.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, pkg.ErrConflict):
		status, message = http.StatusConflict, err.Error()
	default:
		slog.Error("request failed", "error", err)
	}
	writeJSON(w, status, pkg.ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
