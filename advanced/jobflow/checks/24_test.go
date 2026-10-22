package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"testing"
	"time"
)

func Test24IdentityAndImmutableContract(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	raw := frames(t, []byte("same"), []byte("same"))
	jobs, err := s.Submit(ctx, raw, time.Unix(1000, 0))
	if err != nil || jobs[0].ID == jobs[1].ID || jobs[0].Key != jobs[1].Key {
		t.Fatal("body hash used as task identity", jobs, err)
	}
	f, _ := protocol.Unmarshal(raw[0])
	f.Body = []byte("changed")
	raw[0], _ = protocol.Marshal(f)
	if _, err = s.Submit(ctx, raw, time.Unix(1001, 0)); !errors.Is(err, task.ErrConflict) {
		t.Fatal("immutable submission overwritten", err)
	}
	bad := frames(t, []byte("a"), []byte("b"))
	a, _ := protocol.Unmarshal(bad[0])
	b, _ := protocol.Unmarshal(bad[1])
	copy(b.Header.Ctx[24:], a.Header.Ctx[24:])
	bad[1], _ = protocol.Marshal(b)
	if _, err = s.Submit(ctx, bad, time.Unix(1000, 0)); !errors.Is(err, task.ErrConflict) {
		t.Fatal("local ref collision accepted", err)
	}
}
