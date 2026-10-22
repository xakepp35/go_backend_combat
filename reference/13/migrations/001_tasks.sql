CREATE TABLE IF NOT EXISTS tasks (
 id text PRIMARY KEY CHECK (btrim(id) <> ''),
 title text NOT NULL CHECK (btrim(title) <> ''),
 priority integer NOT NULL CHECK (priority BETWEEN 1 AND 3),
 done boolean NOT NULL DEFAULT false
);
