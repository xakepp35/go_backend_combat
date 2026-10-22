// Demo submits a diamond DAG and waits for its final result.
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"jobflow/internal/protocol"
	"jobflow/internal/task"
	"net/http"
	"os"
	"time"
)

func main() {
	base := os.Getenv("PUB_URL")
	if base == "" {
		base = "http://localhost:8080"
	}
	workflow, _ := task.ParseID(task.NewID())
	var ids [4][16]byte
	for i := range ids {
		ids[i], _ = task.ParseID(task.NewID())
	}
	var stream []byte
	for i := range ids {
		var h protocol.Header
		copy(h.Ctx[:16], workflow[:])
		copy(h.Ctx[16:], ids[i][:])
		copy(h.Dst[:], "echo")
		f := protocol.Frame{Header: h, Body: []byte(fmt.Sprintf("node-%d", i))}
		if i == 1 || i == 2 {
			f.Refs = []uint64{binary.BigEndian.Uint64(ids[0][8:])}
		}
		if i == 3 {
			f.Refs = []uint64{binary.BigEndian.Uint64(ids[1][8:]), binary.BigEndian.Uint64(ids[2][8:])}
		}
		b, err := protocol.Marshal(f)
		if err != nil {
			panic(err)
		}
		stream = append(stream, b...)
	}
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Post(base+"/workflows", "application/octet-stream", bytes.NewReader(stream))
	if err != nil {
		panic(err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 202 {
		panic(string(data))
	}
	fmt.Println("workflow", hex.EncodeToString(workflow[:]), "final task", task.UUIDString(ids[3]))
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		resp, err = c.Get(base + "/tasks/" + task.UUIDString(ids[3]))
		if err != nil {
			panic(err)
		}
		var t task.Task
		err = json.NewDecoder(resp.Body).Decode(&t)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}
		if t.State == task.Succeeded {
			fmt.Println("SUCCEEDED", string(t.Result))
			return
		}
		if t.Terminal() {
			panic(t.Error)
		}
		time.Sleep(250 * time.Millisecond)
	}
	panic("workflow deadline exceeded")
}
