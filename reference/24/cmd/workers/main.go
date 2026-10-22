package main

import (
	"context"
	"fmt"
	"jobflow/internal/pool"
)

func main() {
	input := make(chan string, 3)
	input <- "Купить чай"
	input <- "Записаться к врачу"
	input <- "Позвонить другу"
	close(input)
	if err := pool.Run(context.Background(), 2, input, func(ctx context.Context, title string) { fmt.Println("Обработали:", title) }); err != nil {
		fmt.Println(err)
	}
}
