package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"jobflow/migrations"
	"jobflow/pkg"
)

type Store struct{ pool *pgxpool.Pool }

func Open(ctx context.Context, url string) (*Store, error) {
	if url == "" {
		return nil, errors.New("нужна DATABASE_URL")
	}
	p, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}
	return &Store{pool: p}, nil
}
func (s *Store) Close()                         { s.pool.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }
func translate(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return pkg.ErrNotFound
	}
	var p *pgconn.PgError
	if errors.As(err, &p) {
		switch p.Code {
		case "23505":
			return pkg.ErrConflict
		case "23502", "23503", "23514":
			return fmt.Errorf("%w: %s", pkg.ErrInvalid, p.ConstraintName)
		}
	}
	return err
}
func (s *Store) Migrate(ctx context.Context) error {
	return migrations.Up(ctx, s.pool)
}
func (s *Store) create(ctx context.Context, t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, "INSERT INTO tasks(id,title,priority,done) VALUES($1,$2,$3,$4)", t.ID, t.Title, t.Priority, t.Done)
	return translate(err)
}
func (s *Store) get(ctx context.Context, id string) (pkg.Task, error) {
	var t pkg.Task
	err := s.pool.QueryRow(ctx, "SELECT id,title,priority,done FROM tasks WHERE id=$1", id).Scan(&t.ID, &t.Title, &t.Priority, &t.Done)
	return t, translate(err)
}
func (s *Store) list(ctx context.Context) ([]pkg.Task, error) {
	rows, err := s.pool.Query(ctx, "SELECT id,title,priority,done FROM tasks ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]pkg.Task, 0)
	for rows.Next() {
		var t pkg.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Priority, &t.Done); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
func (s *Store) update(ctx context.Context, t pkg.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, "UPDATE tasks SET title=$2,priority=$3,done=$4 WHERE id=$1", t.ID, t.Title, t.Priority, t.Done)
	if err != nil {
		return translate(err)
	}
	if tag.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}
	return nil
}
func (s *Store) delete(ctx context.Context, id string) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return translate(err)
	}
	if tag.RowsAffected() == 0 {
		return pkg.ErrNotFound
	}
	return nil
}
func (s *Store) complete(ctx context.Context, id string) (pkg.Task, error) {
	var t pkg.Task
	err := s.pool.QueryRow(ctx, "UPDATE tasks SET done=true WHERE id=$1 RETURNING id,title,priority,done", id).Scan(&t.ID, &t.Title, &t.Priority, &t.Done)
	return t, translate(err)
}
