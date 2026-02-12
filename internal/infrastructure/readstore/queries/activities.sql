-- name: ListOffenseActivities :many
SELECT id, kind, created_at, base_id, offense_data, defense_data, scan_data, radar_data, trade_data
FROM game.activities
WHERE base_id = $1 AND kind = 'OFFENSE'
  AND ($2 = '' OR offense_data->>'subtype' = $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: ListDefenseActivities :many
SELECT id, kind, created_at, base_id, offense_data, defense_data, scan_data, radar_data, trade_data
FROM game.activities
WHERE base_id = $1 AND kind = 'DEFENSE'
  AND ($2 = '' OR defense_data->>'subtype' = $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: ListScanActivities :many
SELECT id, kind, created_at, base_id, offense_data, defense_data, scan_data, radar_data, trade_data
FROM game.activities
WHERE base_id = $1 AND kind = 'SCAN'
  AND ($2 = '' OR scan_data->>'subtype' = $2)
ORDER BY created_at DESC
LIMIT $3;

-- name: ListRadarActivities :many
SELECT id, kind, created_at, base_id, offense_data, defense_data, scan_data, radar_data, trade_data
FROM game.activities
WHERE base_id = $1 AND kind = 'RADAR'
ORDER BY created_at DESC
LIMIT $2;

-- name: ListTradeActivities :many
SELECT id, kind, created_at, base_id, offense_data, defense_data, scan_data, radar_data, trade_data
FROM game.activities
WHERE base_id = $1 AND kind = 'TRADE'
ORDER BY created_at DESC
LIMIT $2;
