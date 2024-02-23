-- +goose Up
-- +goose StatementBegin
ALTER TABLE teams DROP COLUMN created_at;
ALTER TABLE teams DROP COLUMN updated_at;
ALTER TABLE teams ADD display_name text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE teams DROP COLUMN display_name;
ALTER TABLE teams ADD created_at timestamptz NULL DEFAULT now();
ALTER TABLE teams ADD updated_at timestamptz NULL DEFAULT now();
-- +goose StatementEnd
