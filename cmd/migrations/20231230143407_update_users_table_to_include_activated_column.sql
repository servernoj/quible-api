-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD activated_at timestamptz NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN activated_at;
-- +goose StatementEnd
