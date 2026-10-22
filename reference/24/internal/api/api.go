package api

import (
	"encoding/json"
	"errors"
	"io"
	"jobflow/pkg"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("write response failed", "error", err)
	}
}
func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := err.Error()
	switch {
	case errors.Is(err, pkg.ErrInvalid):
		status = http.StatusBadRequest
	case errors.Is(err, pkg.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, pkg.ErrConflict):
		status = http.StatusConflict
	}
	if status == http.StatusInternalServerError {
		slog.Error("request failed", "error", err)
		message = "внутренняя ошибка"
	}
	WriteJSON(w, status, pkg.ErrorResponse{Error: message})
}
func Decode(w http.ResponseWriter, r *http.Request, value any) error {
	r.Body = http.MaxBytesReader(w, r.Body, pkg.MaxJSONBodyBytes)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(value); err != nil {
		return pkg.ErrInvalid
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return pkg.ErrInvalid
	}
	return nil
}
func New(service pkg.TaskService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var value pkg.Task
		if err := Decode(w, r, &value); err != nil {
			WriteError(w, err)
			return
		}
		res, err := service.Create(r.Context(), pkg.CreateTaskRequest{
			ID: value.ID, Title: value.Title, Priority: value.Priority, Done: value.Done,
		})
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteJSON(w, http.StatusCreated, res.Task)
	})
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.List(r.Context(), pkg.ListTasksRequest{})
		if err != nil {
			WriteError(w, err)
			return
		}
		// HTTP v1 deliberately keeps an array; the service returns an object.
		WriteJSON(w, http.StatusOK, res.Tasks)
	})
	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.Get(r.Context(), pkg.GetTaskRequest{ID: r.PathValue("id")})
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, res.Task)
	})
	mux.HandleFunc("PUT /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		var value pkg.Task
		if err := Decode(w, r, &value); err != nil {
			WriteError(w, err)
			return
		}
		id := r.PathValue("id")
		if value.ID != "" && value.ID != id {
			WriteError(w, pkg.ErrInvalid)
			return
		}
		value.ID = id
		res, err := service.Update(r.Context(), pkg.UpdateTaskRequest{Task: value})
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, res.Task)
	})
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		_, err := service.Delete(r.Context(), pkg.DeleteTaskRequest{ID: r.PathValue("id")})
		if err != nil {
			WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /tasks/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.Complete(r.Context(), pkg.CompleteTaskRequest{ID: r.PathValue("id")})
		if err != nil {
			WriteError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, res.Task)
	})
	return mux
}
