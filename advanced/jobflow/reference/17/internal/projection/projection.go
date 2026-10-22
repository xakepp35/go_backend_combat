package projection

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"jobflow/internal/postgres"
	"jobflow/internal/task"
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
