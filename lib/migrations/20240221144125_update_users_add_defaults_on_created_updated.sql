-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER created_at SET DEFAULT now();
ALTER TABLE users ALTER updated_at SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER created_at DROP DEFAULT;
ALTER TABLE users ALTER updated_at DROP DEFAULT;
-- +goose StatementEnd
