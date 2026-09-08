-- name: CreateCommandLog :one
INSERT INTO command_logs (machine_id, command_type, issued_by, reason)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListCommandLogsByMachine :many
SELECT * FROM command_logs
WHERE machine_id = $1
ORDER BY executed_at DESC
LIMIT $2;
