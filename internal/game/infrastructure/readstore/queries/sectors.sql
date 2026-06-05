-- Sector scan report queries

-- name: GetScansNear :many
SELECT DISTINCT ON (sector_x, sector_y) *
FROM game.scan_reports
WHERE base_id = $1
  AND ((sector_x - $2) * (sector_x - $2)
    +  (sector_y - $3) * (sector_y - $3))
      <= ($4::int * $4::int)
ORDER BY sector_x, sector_y, is_cloaked ASC, created_at DESC;

-- name: GetScanReportByID :one
SELECT *
FROM game.scan_reports
WHERE id = $1 AND base_id = $2;

-- name: GetLatestScanBefore :one
SELECT *
FROM game.scan_reports
WHERE base_id = $1
  AND sector_x = $2
  AND sector_y = $3
  AND created_at <= $4
ORDER BY is_cloaked ASC, created_at DESC
LIMIT 1;

-- name: GetScanReportByOperationUUID :one
SELECT *
FROM game.scan_reports
WHERE source_type = 'OPERATION' AND source_id = $1;
