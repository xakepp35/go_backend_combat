package postgres

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/zeebo/blake3"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"strings"
	"time"
)

const MaxAttempts = 3

type Event struct {
	Seq    int64     `json:"seq"`
	ID     string    `json:"id"`
	Task   task.Task `json:"task"`
	Prefix []byte    `json:"prefix"`
}
type Result struct {
	Token     string `json:"token"`
	Body      []byte `json:"body"`
	Error     string `json:"error,omitempty"`
	Retryable bool   `json:"retryable"`
}

func stripped(t task.Task) task.Task { t.Body = nil; t.Result = nil; return t }
func digest(b []byte) string         { h := blake3.Sum256(b); return hex.EncodeToString(h[:]) }
func putPayload(ctx context.Context, tx pgx.Tx, body []byte) (string, error) {
	if body == nil {
		body = []byte{}
	} // pgx encodes nil []byte as SQL NULL.
	key := digest(body)
	_, err := tx.Exec(ctx, "INSERT INTO payloads(key,raw) VALUES($1,$2) ON CONFLICT DO NOTHING", key, body)
	return key, err
}
func saveJob(ctx context.Context, tx pgx.Tx, t task.Task) error {
	doc, err := json.Marshal(stripped(t))
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE jobs SET doc=$2 WHERE id=$1", t.ID, doc)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO events(id,task_id,version,doc,prefix) SELECT $2,id,$3,$4,prefix FROM jobs WHERE id=$1", t.ID, task.NewID(), t.Version, doc)
	return err
}
func loadJob(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, id string, lock bool) (task.Task, error) {
	var t task.Task
	var doc []byte
	suffix := ""
	if lock {
		suffix = " FOR UPDATE"
	}
	if err := q.QueryRow(ctx, "SELECT doc FROM jobs WHERE id=$1"+suffix, id).Scan(&doc); err != nil {
		return t, translate(err)
	}
	err := json.Unmarshal(doc, &t)
	return t, err
}
func (s *Store) Job(ctx context.Context, id string) (task.Task, error) {
	t, err := loadJob(ctx, s.Pool, id, false)
	if err != nil {
		return t, err
	}
	var body, result []byte
	err = s.Pool.QueryRow(ctx, "SELECT p.raw,r.raw FROM jobs j LEFT JOIN payloads p ON p.key=j.body_key LEFT JOIN payloads r ON r.key=j.result_key WHERE j.id=$1", id).Scan(&body, &result)
	t.Body = body
	t.Result = result
	return t, err
}
func (s *Store) Jobs(ctx context.Context) ([]task.Task, error) {
	rows, err := s.Pool.Query(ctx, "SELECT doc FROM jobs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []task.Task{}
	for rows.Next() {
		var b []byte
		var t task.Task
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(b, &t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
func (s *Store) Payload(ctx context.Context, key string) ([]byte, error) {
	var b []byte
	err := s.Pool.QueryRow(ctx, "SELECT raw FROM payloads WHERE key=$1", key).Scan(&b)
	return b, translate(err)
}

// Submit accepts one closed workflow; external dependencies must already be submitted.
func (s *Store) Submit(ctx context.Context, raw [][]byte, now time.Time) ([]task.Task, error) {
	if len(raw) == 0 || len(raw) > 128 {
		return nil, task.ErrInvalid
	}
	frames := make([]protocol.Frame, len(raw))
	nodes := make([]task.Task, len(raw))
	locals := map[uint64]string{}
	workflow := ""
	total := 0
	for i, b := range raw {
		total += len(b)
		if total > 8<<20 {
			return nil, task.ErrInvalid
		}
		f, err := protocol.Unmarshal(b)
		if err != nil {
			return nil, err
		}
		frames[i] = f
		var tid [16]byte
		copy(tid[:], f.Header.Ctx[16:])
		id := task.UUIDString(tid)
		if _, err = task.ParseID(id); err != nil {
			return nil, err
		}
		wf := hex.EncodeToString(f.Header.Ctx[:16])
		if wf == strings.Repeat("0", 32) {
			return nil, task.ErrInvalid
		}
		if workflow == "" {
			workflow = wf
		}
		if workflow != wf {
			return nil, task.ErrInvalid
		}
		local := binary.BigEndian.Uint64(tid[8:])
		if _, exists := locals[local]; exists {
			return nil, task.ErrConflict
		}
		locals[local] = id
		kind := string(bytes.TrimRight(f.Header.Dst[:], "\x00"))
		if kind != "echo" && kind != "sum" && kind != "counter" {
			return nil, fmt.Errorf("%w: unknown handler", task.ErrInvalid)
		}
		nodes[i] = task.Task{ID: id, Workflow: wf, Kind: kind, Body: f.Body, NBF: f.Header.NBF, EXP: f.Header.EXP, Quotas: task.Quotas{TX: f.Header.TX, RX: f.Header.RX, FX: f.Header.FX, SX: f.Header.SX}, Key: hex.EncodeToString(f.Header.Key[:]), State: task.Pending, Version: 1, UpdatedAt: now}
	}
	for i, f := range frames {
		for _, ref := range f.Refs {
			id, ok := locals[ref]
			if !ok {
				return nil, fmt.Errorf("%w: unresolved local ref", task.ErrInvalid)
			}
			nodes[i].Refs = append(nodes[i].Refs, id)
		}
		for _, dep := range f.Deps {
			nodes[i].Refs = append(nodes[i].Refs, task.UUIDString(dep))
		}
	}
	if err := task.ValidateDAG(nodes); err != nil {
		return nil, err
	}
	err := s.transaction(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(2026101902)"); err != nil {
			return err
		}
		var count int
		if err := tx.QueryRow(ctx, "SELECT count(*) FROM jobs WHERE workflow=$1", workflow).Scan(&count); err != nil {
			return err
		}
		if count != 0 {
			if count != len(nodes) {
				return task.ErrConflict
			}
			for i, t := range nodes {
				var fp string
				if err := tx.QueryRow(ctx, "SELECT fingerprint FROM jobs WHERE id=$1 AND workflow=$2", t.ID, workflow).Scan(&fp); err != nil {
					return task.ErrConflict
				}
				if fp != digest(raw[i]) {
					return task.ErrConflict
				}
			}
			return nil
		}
		for i, f := range frames {
			for _, dep := range f.Deps {
				id := task.UUIDString(dep)
				if _, err := task.ParseID(id); err != nil {
					return err
				}
				other, err := loadJob(ctx, tx, id, false)
				if err != nil {
					return err
				}
				if other.Workflow == workflow {
					return task.ErrInvalid
				}
			}
			t := nodes[i]
			key, err := putPayload(ctx, tx, t.Body)
			if err != nil {
				return err
			}
			prefix := raw[i][:protocol.HeaderSize+int(f.Header.LinksSize())]
			tid, _ := task.ParseID(t.ID)
			doc, err := json.Marshal(stripped(t))
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, "INSERT INTO jobs(id,workflow,local_ref,fingerprint,prefix,body_key,doc) VALUES($1,$2,$3,$4,$5,$6,$7)", t.ID, workflow, tid[8:], digest(raw[i]), prefix, key, doc)
			if err != nil {
				return err
			}
			if err = saveJob(ctx, tx, t); err != nil {
				return err
			}
			if _, err = tx.Exec(ctx, "DELETE FROM drafts WHERE id=$1", t.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := make([]task.Task, len(nodes))
	for i, t := range nodes {
		out[i], err = s.Job(ctx, t.ID)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func depsDone(ctx context.Context, tx pgx.Tx, t task.Task) (bool, bool, error) {
	done := true
	failed := false
	for _, id := range t.Refs {
		dep, err := loadJob(ctx, tx, id, false)
		if err != nil {
			return false, false, err
		}
		if dep.State != task.Succeeded {
			done = false
		}
		if dep.State == task.Failed || dep.State == task.Expired {
			failed = true
		}
	}
	return done, failed, nil
}
func (s *Store) Claim(ctx context.Context, id string, now time.Time, lease time.Duration) (task.Task, error) {
	var t task.Task
	if lease <= 0 || lease > time.Minute {
		return t, task.ErrInvalid
	}
	err := s.transaction(ctx, func(tx pgx.Tx) error {
		var err error
		t, err = loadJob(ctx, tx, id, true)
		if err != nil {
			return err
		}
		done, _, err := depsDone(ctx, tx, t)
		if err != nil {
			return err
		}
		if err = t.CanStart(now, done); err != nil {
			return err
		}
		t.State = task.Running
		t.Token = task.NewID()
		t.LeaseUntil = now.Add(lease)
		t.Attempts++
		t.Version++
		t.UpdatedAt = now
		return saveJob(ctx, tx, t)
	})
	if err != nil {
		return t, err
	}
	return s.Job(ctx, id)
}
func (s *Store) Complete(ctx context.Context, id string, r Result, now time.Time) (task.Task, error) {
	if len(r.Body) > protocol.MaxFrame || len(r.Error) > 1024 {
		return task.Task{}, task.ErrInvalid
	}
	data, err := json.Marshal(r)
	if err != nil {
		return task.Task{}, err
	}
	receipt := digest(data)
	err = s.transaction(ctx, func(tx pgx.Tx) error {
		t, err := loadJob(ctx, tx, id, true)
		if err != nil {
			return err
		}
		var old string
		err = tx.QueryRow(ctx, "SELECT digest FROM receipts WHERE task_id=$1 AND token=$2", id, r.Token).Scan(&old)
		if err == nil {
			if old != receipt {
				return task.ErrConflict
			}
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if t.State != task.Running || t.Token != r.Token || !now.Before(t.LeaseUntil) || (t.EXP != 0 && now.Unix() >= int64(t.EXP)) {
			return task.ErrStale
		}
		t.Error = r.Error
		t.Version++
		t.UpdatedAt = now
		if r.Error == "" {
			t.State = task.Succeeded
			key, err := putPayload(ctx, tx, r.Body)
			if err != nil {
				return err
			}
			if _, err = tx.Exec(ctx, "UPDATE jobs SET result_key=$2 WHERE id=$1", id, key); err != nil {
				return err
			}
		} else if r.Retryable && t.Attempts < MaxAttempts {
			t.State = task.Pending
			t.AvailableAt = now.Add(time.Second * time.Duration(1<<uint(t.Attempts-1)))
		} else {
			t.State = task.Failed
		}
		if err = saveJob(ctx, tx, t); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO receipts(task_id,token,digest) VALUES($1,$2,$3)", id, r.Token, receipt)
		return err
	})
	if err != nil {
		return task.Task{}, err
	}
	return s.Job(ctx, id)
}
func (s *Store) Effect(ctx context.Context, id, token, slot string, delta int64, now time.Time) (int64, error) {
	var value int64
	if slot == "" || len(slot) > 64 {
		return 0, task.ErrInvalid
	}
	err := s.transaction(ctx, func(tx pgx.Tx) error {
		t, err := loadJob(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if t.State != task.Running || t.Token != token || !now.Before(t.LeaseUntil) || (t.EXP != 0 && now.Unix() >= int64(t.EXP)) {
			return task.ErrStale
		}
		var oldDelta int64
		err = tx.QueryRow(ctx, "SELECT delta,result FROM effects WHERE task_id=$1 AND slot=$2", id, slot).Scan(&oldDelta, &value)
		if err == nil {
			if oldDelta != delta {
				return task.ErrConflict
			}
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		limit := t.Quotas.SX
		if limit == 0 || limit > 4 {
			limit = 4
		}
		var count int
		if err = tx.QueryRow(ctx, "SELECT count(*) FROM effects WHERE task_id=$1", id).Scan(&count); err != nil {
			return err
		}
		if count >= int(limit) {
			return task.ErrQuota
		}
		if err = tx.QueryRow(ctx, "INSERT INTO counters(name,value) VALUES('demo',$1) ON CONFLICT(name) DO UPDATE SET value=counters.value+EXCLUDED.value RETURNING value", delta).Scan(&value); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO effects(task_id,slot,delta,result) VALUES($1,$2,$3,$4)", id, slot, delta, value)
		return err
	})
	return value, err
}

// Maintain expires jobs, bounds abandoned attempts, and propagates dependency failure.
