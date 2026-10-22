// sqlcheck executes a teaching schema and assertions in a transaction that is rolled back.
package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"os"
	"time"
)

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	raw := os.Getenv("JF_TEST_DATABASE_URL")
	if raw == "" {
		return fmt.Errorf("JF_TEST_DATABASE_URL required; SQL checks cannot be skipped")
	}
	conn, err := pgx.Connect(ctx, raw)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	schema := pgx.Identifier{fmt.Sprintf("sql_%d", time.Now().UnixNano())}.Sanitize()
	if _, err = tx.Exec(ctx, "CREATE SCHEMA "+schema+";SET LOCAL search_path TO "+schema); err != nil {
		return err
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, string(data))
	return err
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
