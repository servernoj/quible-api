-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD image bytea NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP image;
-- +goose StatementEnd
