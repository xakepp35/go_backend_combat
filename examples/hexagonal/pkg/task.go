package pkg

import (
	"fmt"
	"strings"
)

type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 2
	PriorityHigh   Priority = 3
)

type Task struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Priority Priority `json:"priority"`
	Done     bool     `json:"done"`
}

func (t Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrInvalid)
	}
	if t.Priority < PriorityLow || t.Priority > PriorityHigh {
		return fmt.Errorf("%w: priority must be 1..3", ErrInvalid)
	}
	return nil
}
