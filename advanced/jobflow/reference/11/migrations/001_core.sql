CREATE TABLE drafts (
    id text PRIMARY KEY,
    doc jsonb NOT NULL
);
CREATE TABLE draft_intents (
    seq bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    draft_id text NOT NULL,
    doc jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
