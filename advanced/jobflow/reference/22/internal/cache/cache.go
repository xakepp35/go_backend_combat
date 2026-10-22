package cache

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/zeebo/blake3"
	"jobflow/internal/task"
	"strings"
	"time"
)

type Source interface {
	Payload(context.Context, string) ([]byte, error)
}
type CAS struct {
	Redis     *redis.Client
	Source    Source
	Namespace string
}

func (c CAS) Get(ctx context.Context, key string) ([]byte, error) {
	cacheKey := c.Namespace + ":cas:" + key
	if c.Redis != nil {
		if b, err := c.Redis.Get(ctx, cacheKey).Bytes(); err == nil {
			h := blake3.Sum256(b)
			if hex.EncodeToString(h[:]) == key {
				return b, nil
			}
			_ = c.Redis.Del(ctx, cacheKey).Err()
		}
	}
	b, err := c.Source.Payload(ctx, key)
	if err != nil {
		return nil, err
	}
	h := blake3.Sum256(b)
	if hex.EncodeToString(h[:]) != key {
		return nil, redis.Nil
	}
	if c.Redis != nil {
		_ = c.Redis.Set(ctx, cacheKey, b, time.Hour).Err()
	}
	return b, nil
}

// Prune also repairs leftovers when Redis was down during an earlier durable GC.
func (c CAS) Prune(ctx context.Context) error {
	if c.Redis == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, next, err := c.Redis.Scan(ctx, cursor, c.Namespace+":cas:*", 100).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			_, err = c.Source.Payload(ctx, strings.TrimPrefix(key, c.Namespace+":cas:"))
			if errors.Is(err, task.ErrNotFound) {
				if err = c.Redis.Del(ctx, key).Err(); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}
