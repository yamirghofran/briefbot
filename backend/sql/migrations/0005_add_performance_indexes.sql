-- +goose Up
-- Add indexes for better read status performance
CREATE INDEX idx_items_is_read ON items(is_read);
CREATE INDEX idx_items_user_id_is_read ON items(user_id, is_read);
CREATE INDEX idx_items_created_at ON items(created_at DESC);

-- +goose Down
-- Remove performance indexes
DROP INDEX IF EXISTS idx_items_is_read;
DROP INDEX IF EXISTS idx_items_user_id_is_read;
DROP INDEX IF EXISTS idx_items_created_at;