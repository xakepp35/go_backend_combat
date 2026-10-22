package api

import (
	"jobflow/internal/app"
	"jobflow/internal/memory"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReferenceHTTPContract(t *testing.T) {
	handler := New(app.New(memory.New()))
	request := func(method, path, body string, want int) string {
		t.Helper()
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		if w.Code != want {
			t.Fatalf("%s %s: %d != %d: %s", method, path, w.Code, want, w.Body.String())
		}
		return w.Body.String()
	}
	request("POST", "/tasks", `{"id":"a","title":"чай","priority":1}`, 201)
	request("POST", "/tasks", `{"id":"a","title":"кофе","priority":2}`, 409)
	if body := request("GET", "/tasks/a", "", 200); !strings.Contains(body, "чай") {
		t.Fatal(body)
	}
	for _, bad := range []string{`{"id":"b","title":" ","priority":1}`, `{"id":"b","title":"x","priority":1,"typo":1}`, `{"id":"b","title":"x","priority":1} {}`, `{`} {
		request("POST", "/tasks", bad, 400)
	}
	request("GET", "/tasks/b", "", 404)
	request("PUT", "/tasks/a", `{"id":"z","title":"x","priority":1}`, 400)
	request("PUT", "/tasks/a", `{"title":"новое","priority":3}`, 200)
	request("POST", "/tasks/a/complete", "", 200)
	request("DELETE", "/tasks/a", "", 204)
	request("GET", "/tasks/a", "", 404)
	if body := request("GET", "/tasks", "", 200); strings.TrimSpace(body) != "[]" {
		t.Fatal("пустой список", body)
	}
}
