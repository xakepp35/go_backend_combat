package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	n := 0
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				n++
			}
		}()
	}
	wg.Wait()
	fmt.Println(n)
}
