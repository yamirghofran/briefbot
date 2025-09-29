-- +goose Up
-- +goose StatementBegin

-- Create podcasts table
CREATE TABLE podcasts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    audio_url TEXT,
    dialogues JSONB,
    duration_seconds INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Create podcast_items junction table for many-to-many relationship
CREATE TABLE podcast_items (
    id SERIAL PRIMARY KEY,
    podcast_id INTEGER REFERENCES podcasts(id) ON DELETE CASCADE,
    item_id INTEGER REFERENCES items(id) ON DELETE CASCADE,
    item_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(podcast_id, item_id)
);

-- Create indexes for better performance
CREATE INDEX idx_podcasts_user_id ON podcasts(user_id);
CREATE INDEX idx_podcasts_status ON podcasts(status);
CREATE INDEX idx_podcasts_created_at ON podcasts(created_at DESC);
CREATE INDEX idx_podcast_items_podcast_id ON podcast_items(podcast_id);
CREATE INDEX idx_podcast_items_item_id ON podcast_items(item_id);
CREATE INDEX idx_podcast_items_order ON podcast_items(podcast_id, item_order);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes first
DROP INDEX IF EXISTS idx_podcast_items_order;
DROP INDEX IF EXISTS idx_podcast_items_item_id;
DROP INDEX IF EXISTS idx_podcast_items_podcast_id;
DROP INDEX IF EXISTS idx_podcasts_created_at;
DROP INDEX IF EXISTS idx_podcasts_status;
DROP INDEX IF EXISTS idx_podcasts_user_id;

-- Drop tables
DROP TABLE IF EXISTS podcast_items;
DROP TABLE IF EXISTS podcasts;

-- +goose StatementEnd
