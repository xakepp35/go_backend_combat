package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	p, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Print(err)
		return
	}
	defer p.Close()
	var title string
	if err = p.QueryRow(ctx, "SELECT title FROM tasks WHERE id=$1", "a").Scan(&title); err != nil {
		log.Print(err)
		return
	}
	fmt.Println(title)
}
