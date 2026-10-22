package api

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// RunHTTP waits for in-flight handlers to drain, with a bounded shutdown deadline.
func RunHTTP(ctx context.Context, addr string, h http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	stopped := make(chan struct{})
	drained := make(chan error, 1)
	go func() {
		select {
		case <-ctx.Done():
			c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := srv.Shutdown(c)
			if err != nil {
				_ = srv.Close()
			}
			drained <- err
		case <-stopped:
		}
	}()
	err := srv.ListenAndServe()
	close(stopped)
	if errors.Is(err, http.ErrServerClosed) {
		return <-drained
	}
	return err
}
