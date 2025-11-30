-- Military operations queries

-- name: GetOperation :one
SELECT id, type, owner_user_id, source_base_id, source_x, source_y, target_x, target_y,
       outbound_depart_at, outbound_arrive_at, return_depart_at, return_arrive_at, completed_at,
       phase, result, units, spy_result, attack_result
FROM military_operations
WHERE id = $1;

-- name: ListOperationsByBase :many
SELECT id, type, owner_user_id, source_base_id, source_x, source_y, target_x, target_y,
       outbound_depart_at, outbound_arrive_at, return_depart_at, return_arrive_at, completed_at,
       phase, result, units, spy_result, attack_result
FROM military_operations
WHERE source_base_id = $1
ORDER BY outbound_depart_at DESC;

-- name: ListActiveOperations :many
SELECT id, type, owner_user_id, source_base_id, source_x, source_y, target_x, target_y,
       outbound_depart_at, outbound_arrive_at, return_depart_at, return_arrive_at, completed_at,
       phase, result, units, spy_result, attack_result
FROM military_operations
WHERE source_base_id = $1 AND phase <> 'COMPLETED'
ORDER BY outbound_depart_at DESC;
