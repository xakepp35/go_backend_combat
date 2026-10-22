package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestProgramProducesOutput(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v: %s", err, out)
	}
	if strings.TrimSpace(string(out)) == "" {
		t.Fatal("ожидали непустой результат")
	}
}
