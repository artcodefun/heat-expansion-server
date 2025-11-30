-- Activities queries

-- name: ListActivitiesByBase :many
SELECT id, kind, created_at, base_id,
       operation_data, scan_data, radar_data, trade_data
FROM activities
WHERE base_id = @base_id
ORDER BY created_at DESC, id DESC
LIMIT $1 OFFSET $2;

-- name: InsertActivity :one
INSERT INTO activities (
    kind, created_at, base_id,
    operation_data, scan_data, radar_data, trade_data
) VALUES (
    @kind, @created_at, @base_id,
    @operation_data, @scan_data, @radar_data, @trade_data
)
RETURNING id;

-- name: DeleteActivitiesByBase :exec
DELETE FROM activities WHERE base_id = @base_id;
