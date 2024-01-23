-- +goose Up
-- +goose StatementBegin
create table chats (
  id uuid primary key default gen_random_uuid (),
  resource text not null,
  summary text null,
  title text not null,
  parent_id uuid null references chats,
  is_private boolean null,
  owner_id uuid null references users
);
create table chat_user (
  id uuid primary key default gen_random_uuid (),
  chat_id uuid not null references chats,
  user_id uuid not null references users,
  is_ro boolean not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_user;
drop table if exists chats;
-- +goose StatementEnd
