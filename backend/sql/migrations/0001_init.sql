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
file_key TEXT,
text_content TEXT,
summary TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS item_types (
id SERIAL PRIMARY KEY,
user_id int,
title TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS authors (
id SERIAL PRIMARY KEY,
user_id int,
title TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS platforms (
id SERIAL PRIMARY KEY,
user_id int,
title TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS tags (
id SERIAL PRIMARY KEY,
user_id INT,
title TEXT,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Junction tables for many-to-many relationships
CREATE TABLE IF NOT EXISTS item_item_types (
item_id INT,
item_type_id INT,
PRIMARY KEY (item_id, item_type_id),
FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
FOREIGN KEY (item_type_id) REFERENCES item_types(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS item_authors (
item_id INT,
author_id INT,
PRIMARY KEY (item_id, author_id),
FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS item_platforms (
item_id INT,
platform_id INT,
PRIMARY KEY (item_id, platform_id),
FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
FOREIGN KEY (platform_id) REFERENCES platforms(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS item_tags (
item_id INT,
tag_id INT,
PRIMARY KEY (item_id, tag_id),
FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- +goose Down
-- Drop tables
DROP TABLE IF EXISTS item_tags;
DROP TABLE IF EXISTS item_platforms;
DROP TABLE IF EXISTS item_authors;
DROP TABLE IF EXISTS item_item_types;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS platforms;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS item_types;
DROP TABLE IF EXISTS items;
