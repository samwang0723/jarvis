-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, phone, password) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $2, last_name = $3, email = $4, phone = $5, password = $6
WHERE id = $1;

-- name: UpdateSessionID :exec
UPDATE users SET session_id = $1, session_expired_at = $2 WHERE id = $3;

-- name: DeleteSessionID :exec
UPDATE users SET session_id = NULL, session_expired_at = NULL WHERE id = $1;

-- name: DeleteUserByID :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
