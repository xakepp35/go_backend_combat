package coursecheck_test

import (
	"os"
	"testing"
)

func Test07MigrationExists(t *testing.T) {
	b, err := os.ReadFile("../migrations/001_core.sql")
	if err != nil || len(b) == 0 {
		t.Fatal("missing first SQL migration", err)
	}
}
