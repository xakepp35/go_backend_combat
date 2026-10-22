package main

import (
	"errors"
	"jobflow/pkg"
	"testing"
)

func TestValidationAndCRUD(t *testing.T) {
	for _, p := range []int{0, 4} {
		if _, err := pkg.New("a", "чай", p); !errors.Is(err, pkg.ErrInvalid) {
			t.Fatal("приоритет", p, err)
		}
	}
	if _, err := pkg.New("a", "   ", 1); !errors.Is(err, pkg.ErrInvalid) {
		t.Fatal("пустое название")
	}
	a, _ := pkg.New("a", "чай", 1)
	m := make(map[string]pkg.Task)
	if err := create(m, a); err != nil {
		t.Fatal(err)
	}
	b := a
	b.Title = "кофе"
	if !errors.Is(create(m, b), pkg.ErrConflict) {
		t.Fatal("дубликат")
	}
	got, _ := get(m, "a")
	if got.Title != "чай" {
		t.Fatal("данные изменены при отказе")
	}
	if !errors.Is(update(m, pkg.Task{ID: "z", Title: "x", Priority: 1}), pkg.ErrNotFound) {
		t.Fatal("неизвестный ID")
	}
	if err := remove(m, "a"); err != nil {
		t.Fatal(err)
	}
}
