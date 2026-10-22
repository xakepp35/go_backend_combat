package main

import "testing"

func TestCRUD(t *testing.T) {
	m := make(map[string]string)
	if !create(m, "a", "чай") || create(m, "a", "кофе") {
		t.Fatal("дубликат")
	}
	if v, ok := get(m, "a"); !ok || v != "чай" {
		t.Fatal("старое значение потеряно")
	}
	if update(m, "z", "x") || !update(m, "a", "кофе") {
		t.Fatal("изменение")
	}
	if !remove(m, "a") || remove(m, "a") {
		t.Fatal("удаление")
	}
}
