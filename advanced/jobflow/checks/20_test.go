package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/postgres"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"testing"
	"time"
)

func Test20TimeAndDependencies(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	raw := dependent(t, frames(t, []byte("a"), []byte("b")))
	f, _ := protocol.Unmarshal(raw[0])
	f.Header.NBF = 1001
	f.Header.EXP = 1003
	raw[0], _ = protocol.Marshal(f)
	ts, err := s.Submit(ctx, raw, time.Unix(1000, 0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Claim(ctx, ts[0].ID, time.Unix(1000, 0), time.Second); !errors.Is(err, task.ErrNotReady) {
		t.Fatal("early start", err)
	}
	if _, err = s.Claim(ctx, ts[1].ID, time.Unix(1001, 0), time.Second); !errors.Is(err, task.ErrNotReady) {
		t.Fatal("dependency skipped", err)
	}
	a, err := s.Claim(ctx, ts[0].ID, time.Unix(1001, 0), 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Complete(ctx, a.ID, postgres.Result{Token: a.Token}, time.Unix(1003, 0)); !errors.Is(err, task.ErrStale) {
		t.Fatal("result at expiration", err)
	}
	if err = s.Maintain(ctx, time.Unix(1003, 0)); err != nil {
		t.Fatal(err)
	}
	a, err = s.Job(ctx, a.ID)
	if err != nil || a.State != task.Expired {
		t.Fatal(a, err)
	}
}

func Test20WaitingBatchDoesNotStarveExpiration(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	for range 2 {
		bodies := make([][]byte, 128)
		raw := frames(t, bodies...)
		for i, b := range raw {
			f, _ := protocol.Unmarshal(b)
			f.Header.NBF = 5000
			raw[i], _ = protocol.Marshal(f)
		}
		if _, err := s.Submit(ctx, raw, time.Unix(1000, 0)); err != nil {
			t.Fatal(err)
		}
	}
	raw := frames(t, []byte("must expire"))
	f, _ := protocol.Unmarshal(raw[0])
	f.Header.EXP = 1001
	raw[0], _ = protocol.Marshal(f)
	ts, err := s.Submit(ctx, raw, time.Unix(1000, 0))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Maintain(ctx, time.Unix(1002, 0)); err != nil {
		t.Fatal(err)
	}
	got, err := s.Job(ctx, ts[0].ID)
	if err != nil || got.State != task.Expired {
		t.Fatal("waiting batch starved expiration", got.State, err)
	}
}
