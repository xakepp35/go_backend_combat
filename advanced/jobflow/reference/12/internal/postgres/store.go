package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"jobflow/internal/task"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Store struct{ Pool *pgxpool.Pool }

func Open(ctx context.Context, dsn string) (*Store, error) {
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return &Store{Pool: p}, nil
}
func (s *Store) Close() { s.Pool.Close() }
func translate(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return task.ErrNotFound
	}
	var pe *pgconn.PgError
	if errors.As(err, &pe) && pe.Code == "23505" {
		return task.ErrConflict
	}
	return err
}

// Migrate is deliberately small: forward-only, checksummed SQL, one transaction per migration.
func (s *Store) Migrate(ctx context.Context, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no migrations in %s", dir)
	}
	sort.Strings(files)
	for _, name := range files {
		data, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		digest := hex.EncodeToString(sum[:])
		version := filepath.Base(name)
		err = s.transaction(ctx, func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(2026101901)"); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations(version text PRIMARY KEY, checksum text NOT NULL)"); err != nil {
				return err
			}
			var old string
			err := tx.QueryRow(ctx, "SELECT checksum FROM schema_migrations WHERE version=$1", version).Scan(&old)
			if err == nil {
				if old != digest {
					return fmt.Errorf("migration %s was changed", version)
				}
				return nil
			}
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
			if _, err = tx.Exec(ctx, string(data)); err != nil {
				return err
			}
			_, err = tx.Exec(ctx, "INSERT INTO schema_migrations VALUES($1,$2)", version, digest)
			return err
		})
		if err != nil {
			return fmt.Errorf("migration %s: %w", version, err)
		}
	}
	return nil
}
func (s *Store) transaction(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	if err = f(tx); err != nil {
		return translate(err)
	}
	return translate(tx.Commit(ctx))
}
func (s *Store) Create(ctx context.Context, t task.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	t.State = task.Draft
	t.Version = 1
	t.UpdatedAt = time.Now().UTC()
	doc, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = s.Pool.Exec(ctx, "INSERT INTO drafts(id,doc) VALUES($1,$2)", t.ID, doc)
	return translate(err)
}
func (s *Store) Get(ctx context.Context, id string) (task.Task, error) {
	var t task.Task
	var data []byte
	err := s.Pool.QueryRow(ctx, "SELECT doc FROM drafts WHERE id=$1", id).Scan(&data)
	if err != nil {
		return t, translate(err)
	}
	err = json.Unmarshal(data, &t)
	return t, err
}
func (s *Store) List(ctx context.Context) ([]task.Task, error) {
	rows, err := s.Pool.Query(ctx, "SELECT doc FROM drafts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []task.Task{}
	for rows.Next() {
		var data []byte
		var t task.Task
		if err = rows.Scan(&data); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
func (s *Store) Update(ctx context.Context, t task.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	return s.transaction(ctx, func(tx pgx.Tx) error {
		var data []byte
		if err := tx.QueryRow(ctx, "SELECT doc FROM drafts WHERE id=$1 FOR UPDATE", t.ID).Scan(&data); err != nil {
			return err
		}
		var old task.Task
		if err := json.Unmarshal(data, &old); err != nil {
			return err
		}
		t.State = task.Draft
		t.Version = old.Version + 1
		t.UpdatedAt = time.Now().UTC()
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "UPDATE drafts SET doc=$2 WHERE id=$1", t.ID, data)
		return err
	})
}
func (s *Store) Delete(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, "DELETE FROM drafts WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return task.ErrNotFound
	}
	return nil
}
