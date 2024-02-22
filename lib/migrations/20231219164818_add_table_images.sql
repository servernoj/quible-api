-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images(
  id TEXT PRIMARY KEY,
  parent_id TEXT NULL REFERENCES images,
  slug TEXT NOT NULL,
  display_name TEXT NOT NULL,
  image_url TEXT NOT NULL
);
CREATE INDEX idx_images_display_name ON images(display_name);
CREATE INDEX idx_images_slug ON images(slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
