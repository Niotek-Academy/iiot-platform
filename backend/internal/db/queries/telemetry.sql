-- name: InsertTelemetry :one
INSERT INTO telemetry_logs (sensor_id, metric_value, recorded_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRecentTelemetryBySensor :many
SELECT * FROM telemetry_logs
WHERE sensor_id = $1
ORDER BY recorded_at DESC
LIMIT $2;
