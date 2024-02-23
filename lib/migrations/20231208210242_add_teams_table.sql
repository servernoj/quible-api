-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
	"name" text NOT NULL,
	arena text NOT NULL,
	color text NOT NULL,
	rsc_id int4 NOT NULL,
	created_at timestamptz NULL DEFAULT now(),
	updated_at timestamptz NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
