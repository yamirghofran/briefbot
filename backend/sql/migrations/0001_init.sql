-- +goose Up
CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
name TEXT,
email TEXT,
auth_provider TEXT,
oauth_id TEXT,
password_hash TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items (
id SERIAL PRIMARY KEY,
user_id int,
url TEXT,
is_read BOOLEAN DEFAULT FALSE,
text_content TEXT,
summary TEXT,
type TEXT,
tags TEXT[],
platform TEXT,
authors TEXT[],
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
-- Drop tables
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;
