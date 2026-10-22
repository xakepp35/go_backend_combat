package main

import (
	"jobflow/pkg"
	"testing"
)

func TestTaskCRUD(t *testing.T) {
	a := pkg.Task{ID: "a", Title: "чай", Priority: 1}
	b := a
	b.Complete()
	if a.Done || !b.Done {
		t.Fatal("копия")
	}
	m := make(map[string]pkg.Task)
	if !create(m, a) || create(m, b) {
		t.Fatal("дубликат")
	}
	if !update(m, b) {
		t.Fatal("изменение")
	}
	got, ok := get(m, "a")
	if !ok || !got.Done {
		t.Fatal("изменённая копия не сохранена")
	}
	if !remove(m, "a") {
		t.Fatal("удаление")
	}
}
