package main

import "jobflow/pkg"

func create(m map[string]pkg.Task, t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := m[t.ID]; ok {
		return pkg.ErrConflict
	}
	m[t.ID] = t
	return nil
}
func get(m map[string]pkg.Task, id string) (pkg.Task, error) {
	t, ok := m[id]
	if !ok {
		return pkg.Task{}, pkg.ErrNotFound
	}
	return t, nil
}
func update(m map[string]pkg.Task, t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	if _, ok := m[t.ID]; !ok {
		return pkg.ErrNotFound
	}
	m[t.ID] = t
	return nil
}
func remove(m map[string]pkg.Task, id string) error {
	if _, ok := m[id]; !ok {
		return pkg.ErrNotFound
	}
	delete(m, id)
	return nil
}
