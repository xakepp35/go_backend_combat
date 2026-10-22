-- +goose Up
CREATE TABLE IF NOT EXISTS task_events (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 task_id text NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
 action text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE task_events;
