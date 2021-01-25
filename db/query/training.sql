-- name: CreateTraining :one
INSERT INTO training (user_id, date, sport, type, intensity, details, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteTraining :exec
UPDATE training
SET deleted_at = now()
WHERE id = $1;

-- name: GetTraining :one
SELECT *
FROM training
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListTrainingsByUserInPeriod :many
SELECT *
FROM training
WHERE user_id = $1
  AND date between $2 AND $3
  AND deleted_at IS NULL
ORDER BY id;

-- name: ListTrainingsByUser :many
SELECT *
FROM training
WHERE user_id = $1
ORDER BY id;

-- name: UpdateTraining :one
UPDATE training
SET date      = $2,
    sport     = $3,
    type      = $4,
    intensity = $5,
    details   = $6,
    status    = $7
WHERE id = $1
RETURNING *;
