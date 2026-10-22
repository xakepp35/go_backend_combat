package api

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"jobflow/internal/app"
	"jobflow/internal/postgres"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHTTPContract(t *testing.T) {
	store := newAPIStore(t)
	handler := New(app.New(store))
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
func newAPIStore(t *testing.T) *postgres.Store {
	t.Helper()
	ctx := context.Background()
	raw := os.Getenv("JF_TEST_DATABASE_URL")
	if raw == "" {
		t.Fatal("нужна реальная база JF_TEST_DATABASE_URL")
	}
	conn, err := pgx.Connect(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("api_%d", time.Now().UnixNano())
	name := pgx.Identifier{schema}.Sanitize()
	if _, err = conn.Exec(ctx, "CREATE SCHEMA "+name); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = conn.Exec(context.Background(), "DROP SCHEMA "+name+" CASCADE")
		_ = conn.Close(context.Background())
	})
	u, _ := url.Parse(raw)
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	s, err := postgres.Open(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	if err = s.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	return s
}
