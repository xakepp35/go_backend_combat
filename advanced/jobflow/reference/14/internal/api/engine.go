package api

import (
	"bufio"
	"errors"
	"io"
	"jobflow/internal/postgres"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func RegisterEngine(m *http.ServeMux, s *postgres.Store) {
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

}
