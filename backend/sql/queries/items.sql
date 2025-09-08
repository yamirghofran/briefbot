-- name: GetItem :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemsByUser :many
SELECT * FROM items WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetUnreadItemsByUser :many
SELECT * FROM items WHERE user_id = $1 AND is_read = FALSE ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (user_id, url, file_key, text_content, summary) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateItem :exec
UPDATE items SET url = $2, is_read = $3, file_key = $4, text_content = $5, summary = $6, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: MarkItemAsRead :exec
UPDATE items SET is_read = TRUE, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;

-- Item Types
-- name: GetItemType :one
SELECT * FROM item_types WHERE id = $1;

-- name: GetItemTypesByUser :many
SELECT * FROM item_types WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateItemType :one
INSERT INTO item_types (user_id, title) VALUES ($1, $2) RETURNING *;

-- name: UpdateItemType :exec
UPDATE item_types SET title = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteItemType :exec
DELETE FROM item_types WHERE id = $1;

-- Authors
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;

-- name: GetAuthorsByUser :many
SELECT * FROM authors WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateAuthor :one
INSERT INTO authors (user_id, title) VALUES ($1, $2) RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors SET title = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;

-- Platforms
-- name: GetPlatform :one
SELECT * FROM platforms WHERE id = $1;

-- name: GetPlatformsByUser :many
SELECT * FROM platforms WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreatePlatform :one
INSERT INTO platforms (user_id, title) VALUES ($1, $2) RETURNING *;

-- name: UpdatePlatform :exec
UPDATE platforms SET title = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeletePlatform :exec
DELETE FROM platforms WHERE id = $1;

-- Tags
-- name: GetTag :one
SELECT * FROM tags WHERE id = $1;

-- name: GetTagsByUser :many
SELECT * FROM tags WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateTag :one
INSERT INTO tags (user_id, title) VALUES ($1, $2) RETURNING *;

-- name: UpdateTag :exec
UPDATE tags SET title = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;
