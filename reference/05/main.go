package main

import (
	"fmt"
)

func create(m map[string]string, id, title string) bool {
	if _, ok := m[id]; ok {
		return false
	}
	m[id] = title
	return true
}
func get(m map[string]string, id string) (string, bool) { v, ok := m[id]; return v, ok }
func update(m map[string]string, id, title string) bool {
	if _, ok := m[id]; !ok {
		return false
	}
	m[id] = title
	return true
}
func remove(m map[string]string, id string) bool {
	if _, ok := m[id]; !ok {
		return false
	}
	delete(m, id)
	return true
}
func main() {
	m := make(map[string]string)
	create(m, "a", "Купить чай")
	v, ok := get(m, "a")
	fmt.Println(v, ok)
	update(m, "a", "Купить зелёный чай")
	fmt.Println(m["a"])
	remove(m, "a")
	_, ok = get(m, "a")
	fmt.Println("После удаления:", ok)
}
