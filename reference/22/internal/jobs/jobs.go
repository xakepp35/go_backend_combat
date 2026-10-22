package jobs

import (
	"context"
	"errors"
	"fmt"
	"jobflow/internal/api"
	"jobflow/pkg"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ErrTemporary = errors.New("временная ошибка обработчика")

type Job struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Status   string `json:"status"`
	Attempts int    `json:"attempts"`
	Result   string `json:"result"`
	Error    string `json:"error"`
}
type Store interface {
	CreateJob(context.Context, Job) (Job, bool, error)
	GetJob(context.Context, string) (Job, error)
	Claim(context.Context) (Job, bool, error)
	Finish(context.Context, Job, string, error) error
	Recover(context.Context) error
}
type Handler func(context.Context, string) (string, error)

func CountWords(ctx context.Context, text string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	result := strconv.Itoa(len(strings.Fields(text)))
	return result, ctx.Err()
}
func safeHandle(ctx context.Context, h Handler, text string) (result string, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("handler panic: %v", p)
		}
	}()
	return h(ctx, text)
}
func Run(ctx context.Context, store Store, workers int, handle Handler) error {
	if workers < 1 || workers > 32 {
		return errors.New("workers must be in 1..32")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	failures := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				job, found, err := store.Claim(ctx)
				if err != nil {
					if ctx.Err() == nil {
						failures <- err
						cancel()
					}
					return
				}
				if !found {
					timer := time.NewTimer(25 * time.Millisecond)
					select {
					case <-ctx.Done():
						timer.Stop()
						return
					case <-timer.C:
					}
					continue
				}
				slog.Info("job started", "job_id", job.ID, "attempt", job.Attempts)
				result, cause := safeHandle(ctx, handle, job.Text)
				cleanup, stop := context.WithTimeout(context.Background(), 2*time.Second)
				err = store.Finish(cleanup, job, result, cause)
				stop()
				if err != nil {
					failures <- err
					cancel()
					return
				}
				slog.Info("job finished", "job_id", job.ID, "attempt", job.Attempts, "error", cause)
			}
		}()
	}
	wg.Wait()
	select {
	case err := <-failures:
		return err
	default:
		return ctx.Err()
	}
}
func Register(mux *http.ServeMux, store Store) {
	mux.HandleFunc("POST /jobs", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		}
		if err := api.Decode(w, r, &input); err != nil {
			api.WriteError(w, err)
			return
		}
		if strings.TrimSpace(input.ID) == "" || strings.TrimSpace(input.Text) == "" {
			api.WriteError(w, pkg.ErrInvalid)
			return
		}
		job, created, err := store.CreateJob(r.Context(), Job{ID: input.ID, Text: input.Text})
		if err != nil {
			api.WriteError(w, err)
			return
		}
		status := 200
		if created {
			status = 202
		}
		api.WriteJSON(w, status, job)
	})
	mux.HandleFunc("GET /jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		job, err := store.GetJob(r.Context(), r.PathValue("id"))
		if err != nil {
			api.WriteError(w, err)
			return
		}
		api.WriteJSON(w, 200, job)
	})
}
