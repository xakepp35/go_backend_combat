// Package migrations owns the schema shipped with this build.
package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

//go:embed *.sql
var files embed.FS

// Up completes pending migrations before the application accepts work.
func Up(ctx context.Context, pool *pgxpool.Pool) error {
	// Goose uses database/sql. This adapter borrows the SAME pgx pool;
	// closing it does not close pool. Application queries stay native pgx.
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	// Never silently adopt a database managed by the old course migrator.
	var legacy bool
	if err := db.QueryRowContext(ctx, "SELECT to_regclass('course_migrations') IS NOT NULL").Scan(&legacy); err != nil {
		return fmt.Errorf("inspect migration history: %w", err)
	}
	if legacy {
		return fmt.Errorf("legacy course_migrations detected: follow docs/POSTGRES_UPGRADE.md before starting this build")
	}
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return fmt.Errorf("migration lock: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, files,
		goose.WithSessionLocker(locker),
		goose.WithDisableGlobalRegistry(true),
	)
	if err != nil {
		return fmt.Errorf("migration provider: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("apply embedded migrations: %w", err)
	}
	return nil
}
