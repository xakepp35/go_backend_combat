package coursecheck_test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"jobflow/internal/postgres"
	"jobflow/internal/projection"
	"jobflow/internal/task"
	"os"
	"testing"
)

func Test17ProjectionOrdering(t *testing.T) {
	r := redis.NewClient(&redis.Options{Addr: os.Getenv("JF_TEST_REDIS_ADDR")})
	defer r.Close()
	p := projection.Projection{Redis: r, Namespace: "test-" + task.NewID()}
	ctx := context.Background()
	id := task.NewID()
	defer r.Del(ctx, p.Key("task:"+id), p.Key("tasks"))
	newer := task.Task{ID: id, Kind: "echo", State: task.Succeeded, Version: 2}
	if err := p.Apply(ctx, postgres.Event{Task: newer}); err != nil {
		t.Fatal(err)
	}
	older := newer
	older.State = task.Pending
	older.Version = 1
	if err := p.Apply(ctx, postgres.Event{Task: older}); err != nil {
		t.Fatal(err)
	}
	got, err := p.Get(ctx, id)
	if err != nil || got.State != task.Succeeded {
		t.Fatal("projection regressed", got, err)
	}
	_ = r.Del(ctx, p.Key("task:"+id)).Err()
	if err = p.Apply(ctx, postgres.Event{Task: newer}); err != nil {
		t.Fatal(err)
	}
	if _, err = p.Get(ctx, id); err != nil {
		t.Fatal("rebuild failed", err)
	}
}
