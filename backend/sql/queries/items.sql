-- name: GetItem :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemsByUser :many
SELECT * FROM items WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetUnreadItemsByUser :many
SELECT * FROM items WHERE user_id = $1 AND is_read = FALSE ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (user_id, url, text_content, summary, type, tags, platform, authors) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateItem :exec
UPDATE items SET url = $2, is_read = $3, text_content = $4, summary = $5, type = $6, tags = $7, platform = $8, authors = $9, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: MarkItemAsRead :exec
UPDATE items SET is_read = TRUE, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;


