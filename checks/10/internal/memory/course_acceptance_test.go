package memory

import (
	"errors"
	"jobflow/pkg"
	"testing"
)

func TestRepositoryContract(t *testing.T) {
	s := New()
	a, _ := pkg.New("a", "чай", 1)
	if err := s.Create(a); err != nil {
		t.Fatal(err)
	}
	b := a
	b.Title = "кофе"
	if !errors.Is(s.Create(b), pkg.ErrConflict) {
		t.Fatal("дубликат")
	}
	got, err := s.Get("a")
	if err != nil || got.Title != "чай" {
		t.Fatal("первое значение потеряно", err)
	}
	if !errors.Is(s.Update(pkg.Task{ID: "z", Title: "x", Priority: 1}), pkg.ErrNotFound) {
		t.Fatal("неизвестное изменение")
	}
	done, err := s.Complete("a")
	if err != nil || !done.Done {
		t.Fatal(err)
	}
	if err = s.Delete("a"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Get("a"); !errors.Is(err, pkg.ErrNotFound) {
		t.Fatal("удаление")
	}
}
