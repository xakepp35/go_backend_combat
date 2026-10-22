package projection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jobflow/internal/postgres"
	"jobflow/internal/task"
	"strconv"
	"strings"
	"time"
)

type Projection struct {
	Redis     *redis.Client
	Namespace string
}

func (p Projection) Key(s string) string { return p.Namespace + ":" + s }

const applyScript = `
local old=tonumber(redis.call('HGET',KEYS[1],'version') or '-1')
if old>=tonumber(ARGV[1]) then return 0 end
redis.call('HSET',KEYS[1],'version',ARGV[1],'doc',ARGV[2])
redis.call('SADD',KEYS[2],ARGV[3])
return 1`

func (p Projection) Apply(ctx context.Context, e postgres.Event) error {
	if e.Task.ID == "" || e.Task.Version < 1 {
		return task.ErrInvalid
	}
	doc, err := json.Marshal(e.Task)
	if err != nil {
		return err
	}
	return p.Redis.Eval(ctx, applyScript, []string{p.Key("task:" + e.Task.ID), p.Key("tasks")}, e.Task.Version, string(doc), e.Task.ID).Err()
}
func (p Projection) Get(ctx context.Context, id string) (task.Task, error) {
	var t task.Task
	b, err := p.Redis.HGet(ctx, p.Key("task:"+id), "doc").Bytes()
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(b, &t)
	return t, err
}
func (p Projection) Ready(ctx context.Context, t task.Task, now time.Time) error {
	done := true
	for _, id := range t.Refs {
		dep, err := p.Get(ctx, id)
		if err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}
		if dep.State != task.Succeeded {
			done = false
		}
	}
	if t.CanStart(now, done) != nil {
		return nil
	}
	// Retrying the signal after a short lease repairs a lost Stream without requiring a new job event.
	return p.Redis.Eval(ctx, `if redis.call('SET',KEYS[1],ARGV[1],'NX','PX',5000) then return redis.call('XADD',KEYS[2],'*','task',ARGV[2],'workflow',ARGV[3]) end return ''`, []string{p.Key(fmt.Sprintf("signal:%s:%d", t.ID, t.Version)), p.Key("ready")}, strconv.FormatInt(t.Version, 10), t.ID, t.Workflow).Err()
}
func (p Projection) Group(ctx context.Context) error {
	err := p.Redis.XGroupCreateMkStream(ctx, p.Key("ready"), "act", "0").Err()
	if err != nil && strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}
func (p Projection) Ack(ctx context.Context, id string) error {
	if err := p.Redis.XAck(ctx, p.Key("ready"), "act", id).Err(); err != nil {
		return err
	}
	return p.Redis.XDel(ctx, p.Key("ready"), id).Err()
}
func (p Projection) Receive(ctx context.Context, consumer string) ([]redis.XMessage, error) {
	claimed, _, err := p.Redis.XAutoClaim(ctx, &redis.XAutoClaimArgs{Stream: p.Key("ready"), Group: "act", Consumer: consumer, MinIdle: 6 * time.Second, Start: "0-0", Count: 1}).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(claimed) > 0 {
		return claimed, nil
	}
	streams, err := p.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{Group: "act", Consumer: consumer, Streams: []string{p.Key("ready"), ">"}, Count: 1, Block: time.Second}).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(streams) == 0 {
		return nil, nil
	}
	return streams[0].Messages, nil
}
