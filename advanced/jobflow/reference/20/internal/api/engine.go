package api

import (
	"bufio"
	"errors"
	"io"
	"jobflow/internal/cache"
	"jobflow/internal/postgres"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func RegisterEngine(m *http.ServeMux, s *postgres.Store, cas cache.CAS, retention time.Duration) {
	m.HandleFunc("POST /workflows", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
		reader := bufio.NewReader(r.Body)
		var frames [][]byte
		for {
			f, err := protocol.Read(reader)
			if err == io.EOF {
				break
			}
			if err != nil {
				Write(w, 400, map[string]string{"error": "invalid frame stream"})
				return
			}
			b, err := protocol.Marshal(f)
			if err != nil {
				Error(w, task.ErrInvalid)
				return
			}
			frames = append(frames, b)
			if len(frames) > 128 {
				Error(w, task.ErrInvalid)
				return
			}
		}
		ts, err := s.Submit(r.Context(), frames, time.Now().UTC())
		if err != nil {
			if errors.Is(err, protocol.ErrFrame) {
				err = task.ErrInvalid
			}
			Error(w, err)
			return
		}
		for _, t := range ts {
			slog.Info("submitted", "workflow", t.Workflow, "task", t.ID)
		}
		Write(w, 202, ts)
	})
	m.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.Job(r.Context(), r.PathValue("id"))
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, t)
	})
	m.HandleFunc("GET /internal/tasks", func(w http.ResponseWriter, r *http.Request) {
		ts, err := s.Jobs(r.Context())
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, ts)
	})
	m.HandleFunc("GET /internal/events", func(w http.ResponseWriter, r *http.Request) {
		after, err := strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
		if err != nil || after < 0 {
			Error(w, task.ErrInvalid)
			return
		}
		es, err := s.Events(r.Context(), after, false)
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, es)
	})
	m.HandleFunc("POST /tasks/{id}/claim", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.Claim(r.Context(), r.PathValue("id"), time.Now().UTC(), 5*time.Second)
		if err != nil {
			Error(w, err)
			return
		}
		slog.Info("claimed", "workflow", t.Workflow, "task", t.ID, "attempt", t.Attempts)
		Write(w, 200, t)
	})
	m.HandleFunc("POST /tasks/{id}/result", func(w http.ResponseWriter, r *http.Request) {
		var result postgres.Result
		if err := Decode(w, r, &result); err != nil {
			Error(w, err)
			return
		}
		t, err := s.Complete(r.Context(), r.PathValue("id"), result, time.Now().UTC())
		if err != nil {
			Error(w, err)
			return
		}
		slog.Info("result", "workflow", t.Workflow, "task", t.ID, "state", t.State)
		Write(w, 200, t)
	})
	m.HandleFunc("POST /tasks/{id}/effects", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
			Slot  string `json:"slot"`
			Delta int64  `json:"delta"`
		}
		if err := Decode(w, r, &req); err != nil {
			Error(w, err)
			return
		}
		n, err := s.Effect(r.Context(), r.PathValue("id"), req.Token, req.Slot, req.Delta, time.Now().UTC())
		if err != nil {
			Error(w, err)
			return
		}
		Write(w, 200, map[string]int64{"value": n})
	})
	m.HandleFunc("GET /payloads/{key}", func(w http.ResponseWriter, r *http.Request) {
		b, err := cas.Get(r.Context(), r.PathValue("key"))
		if err != nil {
			Error(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(b)
	})

}
