package main

import (
	"fmt"
	"jobflow/internal/memory"
	"jobflow/pkg"
)

func main() {
	store := memory.New()
	t, err := pkg.New("a", "Купить чай", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = store.Create(t); err != nil {
		fmt.Println(err)
		return
	}
	saved, err := store.Get("a")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(saved.Title)
}
