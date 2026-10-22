package main

import "jobflow/pkg"

func create(m map[string]pkg.Task, t pkg.Task) bool {
	if _, ok := m[t.ID]; ok {
		return false
	}
	m[t.ID] = t
	return true
}
func get(m map[string]pkg.Task, id string) (pkg.Task, bool) { t, ok := m[id]; return t, ok }
func update(m map[string]pkg.Task, t pkg.Task) bool {
	if _, ok := m[t.ID]; !ok {
		return false
	}
	m[t.ID] = t
	return true
}
func remove(m map[string]pkg.Task, id string) bool {
	if _, ok := m[id]; !ok {
		return false
	}
	delete(m, id)
	return true
}
