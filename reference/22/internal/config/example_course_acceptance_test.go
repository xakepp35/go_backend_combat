package config

import "testing"

func TestReferenceConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("нужна база")
	}
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("WORKER_COUNT", "no")
	if _, err := Load(); err == nil {
		t.Fatal("нужно целое число")
	}
	t.Setenv("WORKER_COUNT", "2")
	c, err := Load()
	if err != nil || c.Workers != 2 {
		t.Fatal(c, err)
	}
}
