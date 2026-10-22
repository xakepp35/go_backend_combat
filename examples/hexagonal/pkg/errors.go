package pkg

import "errors"

var (
	ErrInvalid  = errors.New("invalid task")
	ErrConflict = errors.New("task already exists")
)
