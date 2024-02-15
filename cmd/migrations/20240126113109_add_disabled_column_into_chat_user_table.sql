-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat_user ADD disabled boolean not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat_user DROP COLUMN disabled;
-- +goose StatementEnd
