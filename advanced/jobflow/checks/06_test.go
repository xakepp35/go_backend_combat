package coursecheck_test

import (
	"jobflow/internal/api"
	"jobflow/internal/memory"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test06HTTP(t *testing.T) {
	h := api.New(memory.New())
	for _, c := range []struct {
		method, path, body string
		status             int
	}{{"POST", "/drafts", `{"id":"a","kind":"echo","body":"eA=="}`, 201}, {"POST", "/drafts", `{"id":"a","kind":"echo"}`, 409}, {"POST", "/drafts", `{"id":"bad","kind":""}`, 400}, {"POST", "/drafts", `{} {}`, 400}, {"GET", "/drafts/a", "", 200}, {"GET", "/drafts/bad", "", 404}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(c.method, c.path, strings.NewReader(c.body)))
		if w.Code != c.status {
			t.Fatalf("%s %s = %d %s", c.method, c.path, w.Code, w.Body.String())
		}
	}
}
