-- name: CreateAlert :one
INSERT INTO alerts (machine_id, severity, message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListAlertsByMachine :many
SELECT * FROM alerts
WHERE machine_id = $1
ORDER BY created_at DESC;

-- name: ListAlertsByMachineAndStatus :many
SELECT * FROM alerts
WHERE machine_id = $1
  AND is_resolved = $2
ORDER BY created_at DESC;

-- name: ResolveAlert :one
UPDATE alerts
SET is_resolved = TRUE
WHERE alert_id = $1
RETURNING *;

-- name: CountUnresolvedAlerts :one
SELECT COUNT(*) FROM alerts
WHERE machine_id = $1 AND is_resolved = FALSE;
