-- name: InsertAIEvaluation :one
INSERT INTO ai_evaluations (machine_id, health_score, is_anomaly, rul_hours, evaluated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLatestAIEvaluation :one
SELECT * FROM ai_evaluations
WHERE machine_id = $1
ORDER BY evaluated_at DESC
LIMIT 1;
