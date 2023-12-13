-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  team_id TEXT UNIQUE NOT NULL,
  player_id TEXT NULL,
  image bytea NULL,  
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
