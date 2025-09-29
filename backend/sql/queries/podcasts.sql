-- name: GetPodcast :one
SELECT * FROM podcasts WHERE id = $1;

-- name: GetPodcastByUser :many
SELECT * FROM podcasts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetPodcastsByStatus :many
SELECT * FROM podcasts WHERE status = $1 ORDER BY created_at DESC;

-- name: GetPodcastsByUserAndStatus :many
SELECT * FROM podcasts WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC;

-- name: GetPendingPodcasts :many
SELECT * FROM podcasts 
WHERE status = 'pending' 
ORDER BY created_at ASC 
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: GetProcessingPodcasts :many
SELECT * FROM podcasts WHERE status = 'processing' ORDER BY created_at ASC LIMIT $1;

-- name: GetCompletedPodcasts :many
SELECT * FROM podcasts WHERE status = 'completed' ORDER BY created_at DESC LIMIT $1;

-- name: GetRecentPodcasts :many
SELECT * FROM podcasts WHERE status = 'completed' AND created_at > NOW() - INTERVAL '7 days' ORDER BY created_at DESC LIMIT $1;

-- name: CreatePodcast :one
INSERT INTO podcasts (user_id, title, description, status) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreatePodcastWithDialogues :one
INSERT INTO podcasts (user_id, title, description, status, dialogues) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdatePodcast :exec
UPDATE podcasts SET title = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdatePodcastStatus :exec
UPDATE podcasts SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdatePodcastsStatus :exec
UPDATE podcasts SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = ANY($1::int[]);

-- name: UpdatePodcastStatusWithAudio :exec
UPDATE podcasts SET status = $2, audio_url = $3, duration_seconds = $4, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdatePodcastDialogues :exec
UPDATE podcasts SET dialogues = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdatePodcastAudio :exec
UPDATE podcasts SET audio_url = $2, duration_seconds = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeletePodcast :exec
DELETE FROM podcasts WHERE id = $1;

-- name: GetPodcastItems :many
SELECT items.*, podcast_items.item_order 
FROM items 
JOIN podcast_items ON items.id = podcast_items.item_id 
WHERE podcast_items.podcast_id = $1 
ORDER BY podcast_items.item_order ASC;

-- name: GetPodcastItemIDs :many
SELECT item_id FROM podcast_items WHERE podcast_id = $1 ORDER BY item_order ASC;

-- name: AddItemToPodcast :one
INSERT INTO podcast_items (podcast_id, item_id, item_order) VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveItemFromPodcast :exec
DELETE FROM podcast_items WHERE podcast_id = $1 AND item_id = $2;

-- name: UpdatePodcastItemOrder :exec
UPDATE podcast_items SET item_order = $3 WHERE podcast_id = $1 AND item_id = $2;

-- name: ClearPodcastItems :exec
DELETE FROM podcast_items WHERE podcast_id = $1;

-- name: CountPodcastItems :one
SELECT COUNT(*) FROM podcast_items WHERE podcast_id = $1;

-- name: GetPodcastsForItem :many
SELECT podcasts.* 
FROM podcasts 
JOIN podcast_items ON podcasts.id = podcast_items.podcast_id 
WHERE podcast_items.item_id = $1 
ORDER BY podcasts.created_at DESC;

-- name: GetUserPodcastStats :one
SELECT 
    COUNT(*) as total_podcasts,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_podcasts,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_podcasts,
    COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_podcasts,
    COUNT(CASE WHEN status = 'writing' THEN 1 END) as writing_podcasts,
    COUNT(CASE WHEN created_at > NOW() - INTERVAL '7 days' THEN 1 END) as recent_podcasts
FROM podcasts 
WHERE user_id = $1;

-- name: GetPodcastWithItems :one
SELECT 
    podcasts.*,
    COALESCE(jsonb_agg(
        jsonb_build_object(
            'id', items.id,
            'title', items.title,
            'url', items.url,
            'summary', items.summary,
            'item_order', podcast_items.item_order
        ) ORDER BY podcast_items.item_order
    ) FILTER (WHERE items.id IS NOT NULL), '[]'::jsonb) as items
FROM podcasts
LEFT JOIN podcast_items ON podcasts.id = podcast_items.podcast_id
LEFT JOIN items ON podcast_items.item_id = items.id
WHERE podcasts.id = $1
GROUP BY podcasts.id;