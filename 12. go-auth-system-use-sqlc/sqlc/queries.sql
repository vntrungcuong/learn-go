-- name: CreateUser :one
INSERT INTO users (email, hashed_password, fullname)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3, fullname = $4, updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdatePassword :exec
UPDATE users SET hashed_password = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
