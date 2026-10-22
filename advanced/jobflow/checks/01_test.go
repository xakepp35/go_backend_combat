package coursecheck_test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test01Hello(t *testing.T) {
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "run", "./cmd/hello")
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "JobFlow") {
		t.Fatalf("Hello: %s %v", out, err)
	}
}
