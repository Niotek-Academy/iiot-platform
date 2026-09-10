-- name: CreateMachine :one
INSERT INTO machines (machine_id, name, location, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMachine :one
SELECT * FROM machines
WHERE machine_id = $1;

-- name: ListMachines :many
SELECT * FROM machines
ORDER BY created_at ASC;

-- name: UpdateMachine :one
UPDATE machines
SET name = COALESCE(sqlc.narg('name'), name),
    location = COALESCE(sqlc.narg('location'), location),
    status = COALESCE(sqlc.narg('status'), status)
WHERE machine_id = sqlc.arg('machine_id')
RETURNING *;

-- name: DeleteMachine :exec
DELETE FROM machines
WHERE machine_id = $1;
