package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jobflow/internal/client"
	"jobflow/internal/pool"
	"jobflow/internal/postgres"
	"jobflow/internal/projection"
	"jobflow/internal/task"
	"log/slog"
	"strconv"
	"time"
)

type Effects interface {
	Effect(context.Context, task.Task, string, int64) (int64, error)
}
type Budget struct {
	TX, RX, FX, SX uint64
	limits         [4]uint64
}

func NewBudget(q task.Quotas) *Budget {
	limits := [4]uint64{1 << 20, 1 << 20, 100000, 4}
	for i, v := range []uint32{q.TX, q.RX, q.FX, q.SX} {
		if v != 0 && uint64(v) < limits[i] {
			limits[i] = uint64(v)
		}
	}
	return &Budget{limits: limits}
}
func (b *Budget) Consume(tx, rx, fx, sx uint64) error {
	values := []uint64{tx, rx, fx, sx}
	old := []uint64{b.TX, b.RX, b.FX, b.SX}
	for i, v := range values {
		if v > b.limits[i]-old[i] {
			return task.ErrQuota
		}
	}
	b.TX += tx
	b.RX += rx
	b.FX += fx
	b.SX += sx
	return nil
}

type Handler func(context.Context, task.Task, *Budget, Effects) ([]byte, error)

func Registry() map[string]Handler {
	return map[string]Handler{
		"echo": func(ctx context.Context, t task.Task, b *Budget, e Effects) ([]byte, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if err := b.Consume(0, 0, uint64(len(t.Body)), 0); err != nil {
				return nil, err
			}
			return append([]byte(nil), t.Body...), nil
		},
		"sum": func(ctx context.Context, t task.Task, b *Budget, e Effects) ([]byte, error) {
			var values []int64
			if err := json.Unmarshal(t.Body, &values); err != nil {
				return nil, task.ErrInvalid
			}
			var sum int64
			for _, v := range values {
				if err := ctx.Err(); err != nil {
					return nil, err
				}
				if err := b.Consume(0, 0, 1, 0); err != nil {
					return nil, err
				}
				next := sum + v
				if (v > 0 && next < sum) || (v < 0 && next > sum) {
					return nil, task.ErrInvalid
				}
				sum = next
			}
			return []byte(strconv.FormatInt(sum, 10)), nil
		},
		"counter": func(ctx context.Context, t task.Task, b *Budget, e Effects) ([]byte, error) {
			var delta int64
			if err := json.Unmarshal(t.Body, &delta); err != nil {
				return nil, task.ErrInvalid
			}
			request, err := json.Marshal(struct {
				Token string `json:"token"`
				Slot  string `json:"slot"`
				Delta int64  `json:"delta"`
			}{t.Token, "counter", delta})
			if err != nil {
				return nil, err
			}
			if err := b.Consume(uint64(len(request)), 32, 1, 1); err != nil {
				return nil, err
			}
			value, err := e.Effect(ctx, t, "counter", delta)
			if err != nil {
				return nil, err
			}
			return []byte(strconv.FormatInt(value, 10)), nil
		},
	}
}
func Execute(ctx context.Context, t task.Task, registry map[string]Handler, e Effects) (body []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			body = nil
			err = fmt.Errorf("handler panic: %v", p)
		}
	}()
	h, ok := registry[t.Kind]
	if !ok {
		return nil, task.ErrInvalid
	}
	return h(ctx, t, NewBudget(t.Quotas), e)
}
func Retryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, task.ErrInvalid) || errors.Is(err, task.ErrQuota) {
		return false
	}
	var h *client.HTTPError
	if errors.As(err, &h) {
		return h.Status >= 500 || h.Status == 408
	}
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
func Run(ctx context.Context, p projection.Projection, c *client.Client, workers int) error {
	if workers < 1 || workers > 32 {
		return task.ErrInvalid
	}
	jobs := make(chan func(context.Context), workers)
	for i := range workers {
		jobs <- func(ctx context.Context) {
			consumer := fmt.Sprintf("%s-%d", task.NewID(), i)
			for ctx.Err() == nil {
				if err := p.Group(ctx); err != nil {
					wait(ctx, time.Second)
					continue
				}
				msgs, err := p.Receive(ctx, consumer)
				if err != nil {
					wait(ctx, time.Second)
					continue
				}
				for _, m := range msgs {
					id, ok := m.Values["task"].(string)
					if !ok {
						_ = p.Ack(ctx, m.ID)
						continue
					}
					t, err := c.Claim(ctx, id)
					if err != nil {
						var he *client.HTTPError
						if errors.As(err, &he) && he.Status < 500 {
							_ = p.Ack(ctx, m.ID)
						}
						continue
					}
					slog.Info("executing", "workflow", t.Workflow, "task", t.ID, "attempt", t.Attempts)
					execCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
					body, execErr := Execute(execCtx, t, Registry(), c)
					cancel()
					result := postgres.Result{Token: t.Token, Body: body, Retryable: Retryable(execErr)}
					if execErr != nil {
						result.Error = execErr.Error()
					}
					for attempt := 0; attempt < 3; attempt++ {
						err = c.Complete(ctx, t.ID, result)
						if err == nil {
							break
						}
						var he *client.HTTPError
						if errors.As(err, &he) && he.Status < 500 {
							break
						}
						wait(ctx, time.Duration(attempt+1)*100*time.Millisecond)
					}
					if err == nil {
						_ = p.Ack(ctx, m.ID)
					}
				}
			}
		}
	}
	close(jobs)
	return pool.Run(ctx, workers, jobs)
}
func wait(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
