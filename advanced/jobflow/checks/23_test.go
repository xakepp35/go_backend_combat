package coursecheck_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jobflow/internal/task"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func Test23FourProcesses(t *testing.T) {
	_, dsn := database(t)
	brokers, topic := kafkaTopic(t)
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	bins := t.TempDir()
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", bins+string(filepath.Separator), "./cmd/pub", "./cmd/sub", "./cmd/act", "./cmd/ret")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build services: %s %v", out, err)
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	_ = l.Close()
	base := "http://" + addr
	env := append(os.Environ(), "DATABASE_URL="+dsn, "REDIS_ADDR="+os.Getenv("JF_TEST_REDIS_ADDR"), "KAFKA_BROKERS="+strings.Join(brokers, ","), "KAFKA_TOPIC="+topic, "KAFKA_GROUP=group-"+task.NewID(), "JF_NAMESPACE=e2e-"+task.NewID(), "PUB_URL="+base, "HTTP_ADDR="+addr, "WORKER_COUNT=2")
	for _, svc := range []string{"pub", "sub", "act", "ret"} {
		logPath := filepath.Join(bins, svc+".log")
		log, err := os.Create(logPath)
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(filepath.Join(bins, svc))
		cmd.Dir = root
		cmd.Env = env
		cmd.Stdout = log
		cmd.Stderr = log
		if err = cmd.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = cmd.Process.Signal(os.Interrupt)
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				_ = cmd.Process.Kill()
				<-done
			}
			_ = log.Close()
			if t.Failed() {
				data, _ := os.ReadFile(logPath)
				if len(data) > 4000 {
					data = data[len(data)-4000:]
				}
				t.Logf("%s: %s", svc, data)
			}
		})
	}
	c := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(10 * time.Second)
	ready := false
	for time.Now().Before(deadline) {
		resp, err := c.Get(base + "/healthz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == 200 {
				ready = true
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !ready {
		t.Fatal("pub not ready")
	}
	raw := dependent(t, frames(t, []byte("first"), []byte("last")))
	resp, err := c.Post(base+"/workflows", "application/octet-stream", bytes.NewReader(bytes.Join(raw, nil)))
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != 202 {
		t.Fatalf("submit %d: %s", resp.StatusCode, data)
	}
	var jobs []task.Task
	if err = json.Unmarshal(data, &jobs); err != nil {
		t.Fatal(err)
	}
	deadline = time.Now().Add(35 * time.Second)
	for time.Now().Before(deadline) {
		resp, err = c.Get(base + "/tasks/" + jobs[1].ID)
		if err != nil {
			t.Fatal(err)
		}
		var got task.Task
		err = json.NewDecoder(resp.Body).Decode(&got)
		_ = resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if got.State == task.Succeeded {
			if string(got.Result) != "last" {
				t.Fatal("result changed")
			}
			return
		}
		if got.Terminal() {
			t.Fatal(fmt.Sprintf("workflow failed: %+v", got))
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("four-service workflow timed out")
}
