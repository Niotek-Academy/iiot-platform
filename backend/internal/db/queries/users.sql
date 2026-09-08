-- name: CreateUser :one
INSERT INTO users (username, password_hash, role)
VALUES ($1, $2, COALESCE($3, 'OPERATOR'))
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;