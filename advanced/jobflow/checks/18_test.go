package coursecheck_test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"jobflow/internal/projection"
	"jobflow/internal/task"
	"os"
	"testing"
	"time"
)

func Test18PendingRecovery(t *testing.T) {
	r := redis.NewClient(&redis.Options{Addr: os.Getenv("JF_TEST_REDIS_ADDR")})
	defer r.Close()
	p := projection.Projection{Redis: r, Namespace: "test-" + task.NewID()}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	job := task.Task{ID: task.NewID(), Kind: "echo", State: task.Pending, Version: 1}
	defer r.Del(context.Background(), p.Key("ready"), p.Key("signal:"+job.ID+":1"))
	if err := p.Group(ctx); err != nil {
		t.Fatal(err)
	}
	if err := p.Ready(ctx, job, time.Now()); err != nil {
		t.Fatal(err)
	}
	msgs, err := p.Receive(ctx, "dead-consumer")
	if err != nil || len(msgs) != 1 {
		t.Fatal("receive", msgs, err)
	}
	claimed, _, err := r.XAutoClaim(ctx, &redis.XAutoClaimArgs{Stream: p.Key("ready"), Group: "act", Consumer: "replacement", MinIdle: 0, Start: "0-0", Count: 1}).Result()
	if err != nil || len(claimed) != 1 || claimed[0].ID != msgs[0].ID {
		t.Fatal("lost pending", claimed, err)
	}
	if err = p.Ack(ctx, claimed[0].ID); err != nil {
		t.Fatal(err)
	}
}
