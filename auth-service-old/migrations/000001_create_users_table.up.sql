CREATE TABLE IF NOT EXISTS users(
  id serial PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  hashed_password TEXT NOT NULL,
  full_name TEXT NULL,
  phone TEXT NULL,
  image TEXT NULL,
  is_oauth boolean NOT NULL DEFAULT false,	
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
