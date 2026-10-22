package task

import (
	"errors"
	"time"
)

type State string

const Draft State = "DRAFT"

var ErrInvalid = errors.New("invalid task")
var ErrNotFound = errors.New("task not found")
var ErrConflict = errors.New("task conflict")
var ErrNotReady = errors.New("task not ready")
var ErrExpired = errors.New("task expired")
var ErrStale = errors.New("stale attempt")
var ErrQuota = errors.New("quota exceeded")

type Task struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Body      []byte    `json:"body"`
	State     State     `json:"state"`
	Version   int64     `json:"version"`
	NBF       uint32    `json:"nbf,omitempty"`
	EXP       uint32    `json:"exp,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(id, kind string, body []byte) Task {
	return Task{ID: id, Kind: kind, Body: append([]byte(nil), body...), State: Draft, Version: 1}
}
func Clone(t Task) Task        { t.Body = append([]byte(nil), t.Body...); return t }
func (t Task) Validate() error { return nil }
