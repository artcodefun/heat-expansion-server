-- Activity feed queries

-- name: ListActivities :many
SELECT id, kind, created_at, base_id, operation_data, scan_data, radar_data, trade_data
FROM activities
WHERE base_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: ListActivitiesByKind :many
SELECT id, kind, created_at, base_id, operation_data, scan_data, radar_data, trade_data
FROM activities
WHERE base_id = $1 AND kind = $2
ORDER BY created_at DESC
LIMIT $3;