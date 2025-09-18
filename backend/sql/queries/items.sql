-- name: GetItem :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemsByUser :many
SELECT * FROM items WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetUnreadItemsByUser :many
SELECT * FROM items WHERE user_id = $1 AND is_read = FALSE ORDER BY created_at DESC;

-- name: CreateItem :one
INSERT INTO items (user_id, title, url, text_content, summary, type, tags, platform, authors, processing_status, processing_error) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: CreatePendingItem :one
INSERT INTO items (user_id, title, url, processing_status) VALUES ($1, $2, $3, 'pending') RETURNING *;

-- name: UpdateItem :exec
UPDATE items SET title = $2, url = $3, is_read = $4, text_content = $5, summary = $6, type = $7, tags = $8, platform = $9, authors = $10, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdateItemProcessingStatus :exec
UPDATE items SET processing_status = $2, processing_error = $3, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: MarkItemAsRead :exec
UPDATE items SET is_read = TRUE, modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items WHERE id = $1;

-- name: GetPendingItems :many
SELECT * FROM items WHERE processing_status = 'pending' ORDER BY created_at ASC LIMIT $1;

-- name: UpdateItemAsProcessing :exec
UPDATE items SET processing_status = 'processing', modified_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: GetItemsByProcessingStatus :many
SELECT * FROM items WHERE processing_status = $1 ORDER BY created_at DESC;

-- name: GetFailedItemsForRetry :many
SELECT * FROM items WHERE processing_status = 'failed' AND created_at > NOW() - INTERVAL '24 hours' ORDER BY created_at ASC LIMIT $1;

-- name: GetUnreadItemsFromPreviousDay :many
SELECT * FROM items 
WHERE created_at >= DATE_TRUNC('day', NOW() - INTERVAL '1 day') 
  AND created_at < DATE_TRUNC('day', NOW())
  AND is_read = FALSE
  AND processing_status = 'completed'
ORDER BY created_at DESC;

-- name: GetUnreadItemsFromPreviousDayByUser :many
SELECT * FROM items 
WHERE user_id = $1
  AND created_at >= DATE_TRUNC('day', NOW() - INTERVAL '1 day') 
  AND created_at < DATE_TRUNC('day', NOW())
  AND is_read = FALSE
  AND processing_status = 'completed'
ORDER BY created_at DESC;


