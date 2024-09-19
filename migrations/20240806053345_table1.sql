-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS birthday (name TEXT, id TEXT, Date TEXT, Subscribers JSONB);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE birthday;
SELECT 'down SQL query';
-- +goose StatementEnd
