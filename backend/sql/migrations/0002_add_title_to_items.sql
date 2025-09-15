-- +goose Up
ALTER TABLE items ADD COLUMN title TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE items DROP COLUMN title;
