-- name: CreateSensor :one
INSERT INTO sensors (sensor_id, machine_id, metric_name, unit)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSensor :one
SELECT * FROM sensors
WHERE sensor_id = $1;

-- name: ListSensorsByMachine :many
SELECT * FROM sensors
WHERE machine_id = $1
ORDER BY created_at ASC;

-- name: DeleteSensor :exec
DELETE FROM sensors
WHERE sensor_id = $1;
