// Package task contains domain rules. It does not import database or transport clients.
package task

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("invalid task")
	ErrNotFound = errors.New("task not found")
	ErrConflict = errors.New("task conflict")
	ErrNotReady = errors.New("task not ready")
	ErrExpired  = errors.New("task expired")
	ErrStale    = errors.New("stale attempt")
	ErrQuota    = errors.New("quota exceeded")
)

type State string

const (
	Draft     State = "DRAFT"
	Pending   State = "PENDING"
	Running   State = "RUNNING"
	Succeeded State = "SUCCEEDED"
	Failed    State = "FAILED"
	Expired   State = "EXPIRED"
)

type Quotas struct {
	TX uint32 `json:"tx"`
	RX uint32 `json:"rx"`
	FX uint32 `json:"fx"`
	SX uint32 `json:"sx"`
}

type Task struct {
	ID          string    `json:"id"`
	Workflow    string    `json:"workflow,omitempty"`
	Kind        string    `json:"kind"`
	Body        []byte    `json:"body"`
	Refs        []string  `json:"refs,omitempty"`
	NBF         uint32    `json:"nbf,omitempty"`
	EXP         uint32    `json:"exp,omitempty"`
	Quotas      Quotas    `json:"quotas"`
	State       State     `json:"state"`
	Version     int64     `json:"version"`
	Key         string    `json:"key,omitempty"`
	Result      []byte    `json:"result,omitempty"`
	Attempts    int       `json:"attempts"`
	Token       string    `json:"token,omitempty"`
	LeaseUntil  time.Time `json:"lease_until"`
	AvailableAt time.Time `json:"available_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Error       string    `json:"error,omitempty"`
}

func New(id, kind string, body []byte) Task {
	return Task{ID: id, Kind: kind, Body: append([]byte(nil), body...), State: Draft, Version: 1}
}

func (t Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.Kind) == "" || len(t.Body) > 1<<20 {
		return fmt.Errorf("%w: id, kind and bounded body required", ErrInvalid)
	}
	if t.EXP != 0 && t.NBF >= t.EXP {
		return fmt.Errorf("%w: nbf must precede exp", ErrInvalid)
	}
	return nil
}
func (t Task) Terminal() bool { return t.State == Succeeded || t.State == Failed || t.State == Expired }
func (t Task) CanStart(now time.Time, dependenciesDone bool) error {
	if t.EXP != 0 && now.Unix() >= int64(t.EXP) {
		return ErrExpired
	}
	if t.State != Pending || now.Unix() < int64(t.NBF) || now.Before(t.AvailableAt) || !dependenciesDone {
		return ErrNotReady
	}
	return nil
}
func Clone(t Task) Task {
	t.Body = append([]byte(nil), t.Body...)
	t.Refs = append([]string(nil), t.Refs...)
	t.Result = append([]byte(nil), t.Result...)
	return t
}

// NewID returns canonical lowercase UUIDv7 without depending on a DB extension.
func NewID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	ms := uint64(time.Now().UnixMilli())
	for i := 5; i >= 0; i-- {
		b[i] = byte(ms)
		ms >>= 8
	}
	b[6] = (b[6] & 15) | 0x70
	b[8] = (b[8] & 63) | 0x80
	return UUIDString(b)
}
func UUIDString(b [16]byte) string {
	s := hex.EncodeToString(b[:])
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}
func ParseID(s string) ([16]byte, error) {
	var b [16]byte
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return b, ErrInvalid
	}
	raw, err := hex.DecodeString(strings.ReplaceAll(s, "-", ""))
	if err != nil || len(raw) != 16 {
		return b, ErrInvalid
	}
	copy(b[:], raw)
	if b[6]>>4 != 7 || b[8]>>6 != 2 {
		return b, fmt.Errorf("%w: UUIDv7 required", ErrInvalid)
	}
	if UUIDString(b) != s {
		return b, ErrInvalid
	}
	return b, nil
}
