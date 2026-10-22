package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	h := Observe(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	a, b := httptest.NewRecorder(), httptest.NewRecorder()
	h.ServeHTTP(a, httptest.NewRequest("GET", "/", nil))
	h.ServeHTTP(b, httptest.NewRequest("GET", "/", nil))
	if a.Code != 204 || len(a.Header().Get("X-Request-ID")) != 32 || a.Header().Get("X-Request-ID") == b.Header().Get("X-Request-ID") {
		t.Fatal(a.Code, a.Header(), b.Header())
	}
}
