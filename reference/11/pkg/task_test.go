package pkg

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	for _, c := range []struct {
		title string
		p     int
		bad   bool
	}{{"Чай", 1, false}, {"Чай", 3, false}, {"", 2, true}, {"  ", 2, true}, {"Чай", 0, true}, {"Чай", 4, true}} {
		_, err := New("a", c.title, c.p)
		if (err != nil) != c.bad {
			t.Fatalf("%+v: %v", c, err)
		}
		if c.bad && !errors.Is(err, ErrInvalid) {
			t.Fatal(err)
		}
	}
}
func TestCopy(t *testing.T) {
	a, _ := New("a", "Чай", 1)
	b := a
	b.Complete()
	if a.Done || !b.Done {
		t.Fatal("изменилось неверное значение")
	}
}
