package postgres

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"jobflow/internal/task"
	"time"
)

func (s *Store) Events(ctx context.Context, after int64, pending bool) ([]Event, error) {
	query := "SELECT seq,id,doc,prefix FROM events WHERE seq>$1"
	if pending {
		query += " AND NOT delivered"
	}
	query += " ORDER BY seq LIMIT 500"
	rows, err := s.Pool.Query(ctx, query, after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Event{}
	for rows.Next() {
		var e Event
		var b []byte
		if err = rows.Scan(&e.Seq, &e.ID, &b, &e.Prefix); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(b, &e.Task); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (s *Store) Delivered(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, "UPDATE events SET delivered=true WHERE id=$1", id)
	return err
}

type ReclaimedTask struct {
	ID       string `json:"id"`
	Workflow string `json:"workflow"`
}
type GCReport struct {
	Removed int64           `json:"removed"`
	Tasks   []ReclaimedTask `json:"tasks"`
}

// GC removes payload references only after retention, delivery and live dependent checks.
// Job identity, terminal state, receipts, and compact event metadata remain for replay/idempotency.
func (s *Store) GC(ctx context.Context, cutoff time.Time) (int64, error) {
	report, err := s.Collect(ctx, cutoff)
	return report.Removed, err
}

// Collect returns committed task identities for correlation in pub and ret logs.
func (s *Store) Collect(ctx context.Context, cutoff time.Time) (GCReport, error) {
	report := GCReport{Tasks: []ReclaimedTask{}}
	err := s.transaction(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(2026101902)"); err != nil {
			return err
		}
		rows, err := tx.Query(ctx, "SELECT id,doc FROM jobs WHERE body_key IS NOT NULL OR result_key IS NOT NULL FOR UPDATE")
		if err != nil {
			return err
		}
		var candidates []task.Task
		for rows.Next() {
			var id string
			var b []byte
			var t task.Task
			if err = rows.Scan(&id, &b); err != nil {
				rows.Close()
				return err
			}
			if err = json.Unmarshal(b, &t); err != nil {
				rows.Close()
				return err
			}
			if t.Terminal() && t.UpdatedAt.Before(cutoff) {
				candidates = append(candidates, t)
			}
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return err
		}
		for _, t := range candidates {
			var blocked bool
			if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE task_id=$1 AND NOT delivered) OR EXISTS(SELECT 1 FROM jobs WHERE doc->>'state' IN ('PENDING','RUNNING') AND (doc->'refs') ? $1)`, t.ID).Scan(&blocked); err != nil {
				return err
			}
			if blocked {
				continue
			}
			if _, err = tx.Exec(ctx, "UPDATE jobs SET body_key=NULL,result_key=NULL WHERE id=$1", t.ID); err != nil {
				return err
			}
			report.Tasks = append(report.Tasks, ReclaimedTask{ID: t.ID, Workflow: t.Workflow})
		}
		tag, err := tx.Exec(ctx, `DELETE FROM payloads p WHERE NOT EXISTS(SELECT 1 FROM jobs j WHERE j.body_key=p.key OR j.result_key=p.key)`)
		if err != nil {
			return err
		}
		report.Removed = tag.RowsAffected()
		return nil
	})
	if err != nil {
		return GCReport{}, err
	}
	return report, nil
}
