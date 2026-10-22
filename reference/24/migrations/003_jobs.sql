-- +goose Up
CREATE TABLE jobs (
 id text PRIMARY KEY CHECK (btrim(id) <> ''),
 body text NOT NULL CHECK (btrim(body) <> ''),
 status text NOT NULL DEFAULT 'queued' CHECK(status IN ('queued','running','done','failed')),
 attempts integer NOT NULL DEFAULT 0 CHECK(attempts BETWEEN 0 AND 3),
 result text NOT NULL DEFAULT '',
 error text NOT NULL DEFAULT '',
 available_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE jobs;
