-- +goose Up
CREATE INDEX tasks_open_priority_idx ON tasks(priority DESC,id) WHERE done=false;

-- +goose Down
DROP INDEX tasks_open_priority_idx;
