package coursecheck_test

import (
	"context"
	"encoding/hex"
	"github.com/redis/go-redis/v9"
	"github.com/zeebo/blake3"
	"jobflow/internal/cache"
	"jobflow/internal/task"
	"os"
	"testing"
	"time"
)

type bodySource struct{ body []byte }

func (s bodySource) Payload(ctx context.Context, key string) ([]byte, error) { return s.body, nil }
func Test15CAS(t *testing.T) {
	addr := os.Getenv("JF_TEST_REDIS_ADDR")
	if addr == "" {
		t.Fatal("Redis infrastructure required")
	}
	r := redis.NewClient(&redis.Options{Addr: addr, MaxRetries: 0})
	defer r.Close()
	ctx := context.Background()
	ns := "test-" + task.NewID()
	body := []byte{1, 0, 0}
	hash := blake3.Sum256(body)
	key := hex.EncodeToString(hash[:])
	c := cache.CAS{Redis: r, Source: bodySource{body}, Namespace: ns}
	defer r.Del(ctx, ns+":cas:"+key)
	if err := r.Set(ctx, ns+":cas:"+key, "corrupt", time.Minute).Err(); err != nil {
		t.Fatal(err)
	}
	got, err := c.Get(ctx, key)
	if err != nil || string(got) != string(body) {
		t.Fatal("corruption recovery", got, err)
	}
	bad := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", MaxRetries: 0, DialTimeout: time.Millisecond})
	defer bad.Close()
	c.Redis = bad
	got, err = c.Get(ctx, key)
	if err != nil || string(got) != string(body) {
		t.Fatal("Redis loss lost authoritative body", err)
	}
}
