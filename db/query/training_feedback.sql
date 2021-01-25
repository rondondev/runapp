-- name: CreateTrainingFeedback :one
INSERT INTO training_feedback (training_id, borg_scale)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteTrainingFeedback :exec
DELETE
FROM training_feedback
WHERE id = $1;

-- name: GetTrainingFeedback :one
SELECT *
FROM training_feedback
WHERE id = $1
LIMIT 1;

-- name: ListTrainingFeedbacksByUser :many
SELECT tf.*
FROM training_feedback tf
         JOIN training t ON tf.training_id = t.id
WHERE t.user_id = $1
  AND t.deleted_at IS NULL
ORDER BY tf.id;

-- name: ListTrainingFeedbacksByUserInPeriod :many
SELECT tf.*
FROM training_feedback tf
         JOIN training t ON tf.training_id = t.id
WHERE t.user_id = $1
  AND date between $2 AND $3
  AND t.deleted_at IS NULL
ORDER BY tf.id;

-- name: UpdateTrainingFeedback :one
UPDATE training_feedback
SET borg_scale = $2
WHERE id = $1
RETURNING *;
