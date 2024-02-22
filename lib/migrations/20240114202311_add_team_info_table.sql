-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.team_info (
	id INTEGER PRIMARY KEY,
	name text NOT NULL,
	slug text NOT NULL,
	short_name text NOT NULL,
	abbr text NOT NULL,
	arena_name text NOT NULL,
	arena_size INTEGER NOT NULL,
  color text NOT NULL,
  secondary_color text NOT NULL,
  logo text NULL  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_info;
-- +goose StatementEnd
