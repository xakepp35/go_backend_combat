package main

import (
	"fmt"
	"jobflow/pkg"
)

func main() {
	t, err := pkg.New("a", "Купить чай", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t.Title)
	_, err = pkg.New("b", "   ", 2)
	fmt.Println("Ошибка:", err)
}
