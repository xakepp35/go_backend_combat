package coursecheck_test

import (
	"context"
	"errors"
	"jobflow/internal/task"
	"testing"
)

func Test09RollbackAndIntent(t *testing.T) {
	s, _ := database(t)
	ctx := context.Background()
	_, err := s.Pool.Exec(ctx, `CREATE FUNCTION reject_intent() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN IF NEW.draft_id='reject' THEN RAISE EXCEPTION 'injected'; END IF; RETURN NEW; END$$; CREATE TRIGGER rejected BEFORE INSERT ON draft_intents FOR EACH ROW EXECUTE FUNCTION reject_intent()`)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.CreateWithIntent(ctx, task.New("reject", "echo", nil)); err == nil {
		t.Fatal("injection did not fail")
	}
	if _, err = s.Get(ctx, "reject"); !errors.Is(err, task.ErrNotFound) {
		t.Fatal("partial transaction committed", err)
	}
	if err = s.CreateWithIntent(ctx, task.New("ok", "echo", nil)); err != nil {
		t.Fatal(err)
	}
	var count int
	if err = s.Pool.QueryRow(ctx, "SELECT count(*) FROM draft_intents").Scan(&count); err != nil || count != 1 {
		t.Fatal("intent count", count, err)
	}
}
