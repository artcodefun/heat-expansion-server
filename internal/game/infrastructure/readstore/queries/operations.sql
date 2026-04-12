-- Military operations queries

-- name: GetOperation :one
SELECT *
FROM game.military_operations
WHERE id = $1;

-- name: GetOperationByUUID :one
SELECT *
FROM game.military_operations
WHERE operation_uuid = $1;

-- name: ListOperationsByBase :many
SELECT *
FROM game.military_operations
WHERE source_base_id = $1
ORDER BY outbound_depart_at DESC;

-- name: ListActiveOperations :many
SELECT *
FROM game.military_operations
WHERE source_base_id = $1 AND phase <> 'COMPLETED'
ORDER BY outbound_depart_at DESC;
