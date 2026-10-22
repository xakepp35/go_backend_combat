package postgres

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"jobflow/internal/task"
	"time"
)

// CreateWithIntent is lesson 09's transactional precursor to the submitted-task Outbox.
// Draft intents are local audit records, not executable jobs or Kafka messages.
func (s *Store) CreateWithIntent(ctx context.Context, t task.Task) error {
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
	return s.transaction(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "INSERT INTO drafts(id,doc) VALUES($1,$2)", t.ID, doc); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, "INSERT INTO draft_intents(draft_id,doc) VALUES($1,$2)", t.ID, doc)
		return err
	})
}
