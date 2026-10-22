package api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReferenceHello(t *testing.T) {
	w := httptest.NewRecorder()
	New().ServeHTTP(w, httptest.NewRequest("GET", "/hello", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "JobFlow") {
		t.Fatal(w.Code, w.Body.String())
	}
}
