-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email, auth_provider, oauth_id, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET name = $2, email = $3, auth_provider = $4, oauth_id = $5, password_hash = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
