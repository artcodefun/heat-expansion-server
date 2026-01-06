-- Sector scan report queries

-- name: GetScansNear :many
SELECT DISTINCT ON (sector_x, sector_y)
       id, base_id, sector_x, sector_y, created_at, type, is_cloaked,
       source_operation_id, name, description, image_url, info
FROM scan_reports
WHERE base_id = $1
  AND ((sector_x - $2) * (sector_x - $2)
    +  (sector_y - $3) * (sector_y - $3))
      <= ($4::int * $4::int)
ORDER BY sector_x, sector_y, created_at DESC;

-- name: GetScanReportByID :one
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked, source_operation_id,
       name, description, image_url, info
FROM scan_reports
WHERE id = $1 AND base_id = $2;

-- name: GetLatestScanBefore :one
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked, source_operation_id,
       name, description, image_url, info
FROM scan_reports
WHERE base_id = $1
  AND sector_x = $2
  AND sector_y = $3
  AND created_at <= $4
ORDER BY created_at DESC
LIMIT 1;

-- name: GetScanReportByOperationID :one
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked, source_operation_id,
       name, description, image_url, info
FROM scan_reports
WHERE source_operation_id = $1;
