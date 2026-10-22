CREATE TABLE payloads (
    key text PRIMARY KEY CHECK (key ~ '^[0-9a-f]{64}$'),
    raw bytea NOT NULL
);
CREATE TABLE jobs (
    id text PRIMARY KEY,
    workflow text NOT NULL,
    local_ref bytea NOT NULL CHECK (octet_length(local_ref)=8),
    fingerprint text NOT NULL,
    prefix bytea NOT NULL CHECK (octet_length(prefix)>=128),
    body_key text REFERENCES payloads(key),
    result_key text REFERENCES payloads(key),
    doc jsonb NOT NULL,
    UNIQUE(workflow, local_ref)
);
CREATE TABLE events (
    seq bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id text UNIQUE NOT NULL,
    task_id text NOT NULL REFERENCES jobs(id),
    version bigint NOT NULL,
    doc jsonb NOT NULL,
    prefix bytea NOT NULL,
    delivered boolean NOT NULL DEFAULT false,
    UNIQUE(task_id, version)
);
CREATE INDEX events_pending ON events(seq) WHERE NOT delivered;
