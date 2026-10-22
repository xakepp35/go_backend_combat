CREATE TABLE receipts (
    task_id text NOT NULL REFERENCES jobs(id),
    token text NOT NULL,
    digest text NOT NULL,
    PRIMARY KEY(task_id,token)
);
CREATE TABLE counters (name text PRIMARY KEY, value bigint NOT NULL);
CREATE TABLE effects (
    task_id text NOT NULL REFERENCES jobs(id),
    slot text NOT NULL,
    delta bigint NOT NULL,
    result bigint NOT NULL,
    PRIMARY KEY(task_id,slot)
);
