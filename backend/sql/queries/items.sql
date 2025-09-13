-- name: GetItem :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemsByUser :many
SELECT * FROM items WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetUnreadItemsByUser :many
SELECT * FROM items WHERE user_id = $1 AND is_read = FALSE ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (user_id, title, url, text_content, summary, type, tags, platform, authors) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateItem :exec
UPDATE items SET title = $2, url = $3, is_read = $4, text_content = $5, summary = $6, type = $7, tags = $8, platform = $9, authors = $10, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: MarkItemAsRead :exec
UPDATE items SET is_read = TRUE, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;


