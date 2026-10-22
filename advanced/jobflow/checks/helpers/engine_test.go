package coursecheck_test

import (
	"context"
	"encoding/binary"
	"jobflow/internal/postgres"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"testing"
	"time"
)

func frames(t *testing.T, bodies ...[]byte) [][]byte {
	t.Helper()
	workflow, _ := task.ParseID(task.NewID())
	var out [][]byte
	for _, body := range bodies {
		id, _ := task.ParseID(task.NewID())
		var h protocol.Header
		copy(h.Ctx[:16], workflow[:])
		copy(h.Ctx[16:], id[:])
		copy(h.Dst[:], "echo")
		b, err := protocol.Marshal(protocol.Frame{Header: h, Body: body})
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, b)
	}
	return out
}
func submit(t *testing.T, s *postgres.Store, bodies ...[]byte) []task.Task {
	t.Helper()
	ts, err := s.Submit(context.Background(), frames(t, bodies...), time.Unix(1000, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	return ts
}
func dependent(t *testing.T, raw [][]byte) [][]byte {
	t.Helper()
	a, _ := protocol.Unmarshal(raw[0])
	b, _ := protocol.Unmarshal(raw[1])
	id := a.Header.Ctx[16:]
	b.Refs = []uint64{binary.BigEndian.Uint64(id[8:])}
	out, err := protocol.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	raw[1] = out
	return raw
}
func delivered(t *testing.T, s *postgres.Store) {
	t.Helper()
	ctx := context.Background()
	for {
		es, err := s.Events(ctx, 0, true)
		if err != nil {
			t.Fatal(err)
		}
		if len(es) == 0 {
			return
		}
		for _, e := range es {
			if err = s.Delivered(ctx, e.ID); err != nil {
				t.Fatal(err)
			}
		}
	}
}
