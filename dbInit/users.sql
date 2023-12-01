CREATE TABLE users (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	username text NOT NULL,
	email text NOT NULL,
	hashed_password text NOT NULL,
	full_name text NOT NULL,
	phone text NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	"refresh" uuid NOT NULL DEFAULT gen_random_uuid(),
	image bytea NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);
CREATE INDEX idx_users_email ON public.users USING btree (email);