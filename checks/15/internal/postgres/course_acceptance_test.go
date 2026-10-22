package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"jobflow/pkg"
	"net/url"
	"os"
	"testing"
	"time"
)

func isolatedStore(t *testing.T) (*Store, *pgx.Conn, string) {
	t.Helper()
	ctx := context.Background()
	raw := os.Getenv("JF_TEST_DATABASE_URL")
	if raw == "" {
		t.Fatal("нужна реальная PostgreSQL: JF_TEST_DATABASE_URL")
	}
	conn, err := pgx.Connect(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("core_%d", time.Now().UnixNano())
	name := pgx.Identifier{schema}.Sanitize()
	if _, err = conn.Exec(ctx, "CREATE SCHEMA "+name); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = conn.Exec(context.Background(), "DROP SCHEMA "+name+" CASCADE")
		_ = conn.Close(context.Background())
	})
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	s, err := Open(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	if err = s.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if err = s.Migrate(ctx); err != nil {
		t.Fatal("повтор миграции", err)
	}
	return s, conn, schema
}
func TestPostgresContract(t *testing.T) {
	ctx := context.Background()
	s, _, schema := isolatedStore(t)
	a, _ := pkg.New("a", "чай", 1)
	if _, err := s.Insert(ctx, pkg.InsertTaskRequest{Task: a}); err != nil {
		t.Fatal(err)
	}
	other := a
	other.Title = "кофе"
	if _, err := s.Insert(ctx, pkg.InsertTaskRequest{Task: other}); !errors.Is(err, pkg.ErrConflict) {
		t.Fatal("дубликат")
	}
	got, err := s.Get(ctx, pkg.GetTaskRequest{ID: "a"})
	if err != nil || got.Task.Title != "чай" {
		t.Fatal(got, err)
	}
	all, err := s.List(ctx, pkg.ListTasksRequest{})
	if err != nil || len(all.Tasks) != 1 {
		t.Fatal(all, err)
	}
	other.ID = "z"
	if _, err := s.Update(ctx, pkg.UpdateTaskRequest{Task: other}); !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("неизвестное изменение")
	}
	done, err := s.Complete(ctx, pkg.CompleteTaskRequest{ID: "a"})
	if err != nil || !done.Task.Done {
		t.Fatal(done, err)
	}
	s.Close()
	u, _ := url.Parse(os.Getenv("JF_TEST_DATABASE_URL"))
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	reopened, err := Open(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	again, err := reopened.Get(ctx, pkg.GetTaskRequest{ID: "a"})
	if err != nil || !again.Task.Done {
		t.Fatal("сохранность после переподключения", again, err)
	}
	if _, err = reopened.Delete(ctx, pkg.DeleteTaskRequest{ID: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err = reopened.Get(ctx, pkg.GetTaskRequest{ID: "a"}); !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("удаление")
	}
}
