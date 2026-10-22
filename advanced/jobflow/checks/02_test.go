package coursecheck_test

import (
	"bytes"
	"jobflow/internal/task"
	"testing"
)

func Test02Model(t *testing.T) {
	body := []byte("work")
	a := task.New("a", "echo", body)
	b := task.New("b", "sum", nil)
	if a.ID != "a" || a.Kind != "echo" || !bytes.Equal(a.Body, body) || b.ID == a.ID {
		t.Fatal("model lost input")
	}
}
