package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/postgres"
	"jobflow/internal/task"
	"testing"
)

func Test08PostgresContract(t *testing.T) {
	s, dsn := database(t)
	ctx := context.Background()
	if err := s.Create(ctx, task.New("a", "echo", []byte("durable"))); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(ctx, task.New("a", "echo", nil)); !errors.Is(err, task.ErrConflict) {
		t.Fatal("duplicate", err)
	}
	other, err := postgres.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	got, err := other.Get(ctx, "a")
	if err != nil || string(got.Body) != "durable" {
		t.Fatal("restart persistence", got, err)
	}
	if err = s.Migrate(ctx, "../migrations"); err != nil {
		t.Fatal("migration rerun", err)
	}
}
