package main

import (
	"fmt"
	"jobflow/pkg"
)

func main() {
	original := pkg.Task{ID: "a", Title: "Купить чай", Priority: 1}
	changed := original
	changed.Complete()
	fmt.Println("Оригинал:", original.Done)
	fmt.Println("Копия:", changed.Done)
	m := map[string]pkg.Task{"a": original}
	value := m["a"]
	value.Complete()
	m["a"] = value
	fmt.Println("Сохранено:", m["a"].Done)
}
