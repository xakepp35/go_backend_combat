package pkg

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalid = errors.New("неверная задача")
var ErrNotFound = errors.New("задача не найдена")
var ErrConflict = errors.New("ID уже занят")

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
	Done     bool   `json:"done"`
}

func (t *Task) Complete() { t.Done = true }
func (t Task) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("%w: нужен ID", ErrInvalid)
	}
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("%w: нужно название", ErrInvalid)
	}
	if t.Priority < 1 || t.Priority > 3 {
		return fmt.Errorf("%w: приоритет от 1 до 3", ErrInvalid)
	}
	return nil
}
func New(id, title string, priority int) (Task, error) {
	t := Task{ID: id, Title: title, Priority: priority}
	if err := t.Validate(); err != nil {
		return Task{}, err
	}
	return t, nil
}
