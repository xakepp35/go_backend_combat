package postgres

import (
	"context"
	"encoding/json"
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
