package main

import (
	"fmt"
	"jobflow/internal/task"
)

func main() { a := task.New("a", "echo", []byte("report")); fmt.Println(a.ID, a.Kind, string(a.Body)) }
