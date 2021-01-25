-- name: CreateUser :one
INSERT INTO users (type, name, email, phone, birth, active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListActiveUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
  AND active = TRUE
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListAllUsers :many
SELECT *
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET type          = $2,
    name          = $3,
    email         = $4,
    phone         = $5,
    birth         = $6,
    active        = $7
WHERE id = $1
RETURNING *;
