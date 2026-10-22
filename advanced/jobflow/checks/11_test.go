package coursecheck_test

import (
	"context"
	"errors"
	"io"
	"jobflow/internal/api"
	"jobflow/internal/pool"
	"net"
	"net/http"
	"testing"
	"time"
)

func Test11Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan func(context.Context))
	done := make(chan error, 1)
	go func() { done <- pool.Run(ctx, 2, jobs) }()
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown waits on empty producer")
	}
	if err := pool.Run(context.Background(), 0, jobs); err == nil {
		t.Fatal("unbounded worker setting accepted")
	}
}

func Test11HTTPDrainsActiveHandler(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()
	_ = listener.Close()
	entered := make(chan struct{})
	release := make(chan struct{})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		_, _ = w.Write([]byte("done"))
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverDone := make(chan error, 1)
	go func() { serverDone <- api.RunHTTP(ctx, addr, handler) }()
	requestDone := make(chan error, 1)
	go func() {
		c := &http.Client{Timeout: 3 * time.Second}
		deadline := time.Now().Add(time.Second)
		for {
			resp, e := c.Get("http://" + addr)
			if e == nil {
				_, e = io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				requestDone <- e
				return
			}
			if time.Now().After(deadline) {
				requestDone <- e
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		close(release)
		t.Fatal("HTTP handler never started")
	}
	cancel()
	var early bool
	select {
	case <-serverDone:
		early = true
	case <-time.After(30 * time.Millisecond):
	}
	close(release)
	if early {
		t.Fatal("RunHTTP returned before active handler drained")
	}
	select {
	case err = <-serverDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown did not finish")
	}
	if err = <-requestDone; err != nil {
		t.Fatal(err)
	}
}
