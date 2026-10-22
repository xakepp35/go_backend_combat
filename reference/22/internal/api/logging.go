package api

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type responseLog struct {
	http.ResponseWriter
	status int
}

func (w *responseLog) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *responseLog) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(p)
}
func Observe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw [16]byte
		if _, err := rand.Read(raw[:]); err != nil {
			http.Error(w, "request ID unavailable", 500)
			return
		}
		id := hex.EncodeToString(raw[:])
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		out := &responseLog{ResponseWriter: w}
		next.ServeHTTP(out, r)
		if out.status == 0 {
			out.status = 200
		}
		slog.Info("request finished", "request_id", id, "method", r.Method, "path", r.URL.Path, "status", out.status, "duration", time.Since(start))
	})
}
