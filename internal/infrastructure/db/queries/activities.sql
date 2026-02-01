-- Activities queries

-- name: ListActivitiesByBase :many
SELECT id, kind, created_at, base_id,
       offense_data, defense_data, scan_data, radar_data, trade_data
FROM activities
WHERE base_id = @base_id
ORDER BY created_at DESC, id DESC
LIMIT $1 OFFSET $2;

-- name: InsertActivity :one
INSERT INTO activities (
    id, kind, created_at, base_id,
    offense_data, defense_data, scan_data, radar_data, trade_data
) VALUES (
    @id, @kind, @created_at, @base_id,
    @offense_data, @defense_data, @scan_data, @radar_data, @trade_data
)
RETURNING id;

-- name: DeleteActivitiesByBase :exec
DELETE FROM activities WHERE base_id = @base_id;

-- name: ExistsForOperation :one
SELECT EXISTS (
    SELECT 1 FROM activities
    WHERE base_id = sqlc.arg('base_id') AND kind = sqlc.arg('kind')
      AND (
          (kind = 'OFFENSE' AND (offense_data->>'op_id')::bigint = sqlc.arg('op_id')::bigint) OR
          (kind = 'DEFENSE' AND (defense_data->>'op_id')::bigint = sqlc.arg('op_id')::bigint) OR
          (kind = 'RADAR' AND (radar_data->>'op_id')::bigint = sqlc.arg('op_id')::bigint)
      )
);

-- name: ExistsForScanReport :one
SELECT EXISTS (
    SELECT 1 FROM activities
    WHERE kind = 'SCAN' AND (scan_data->>'report_id')::bigint = sqlc.arg('report_id')::bigint
);
