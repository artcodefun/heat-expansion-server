-- Scan report queries

-- name: InsertScanReport :one
INSERT INTO game.scan_reports (
    base_id, sector_x, sector_y, created_at, type, is_cloaked,
    source_operation_id, source_scanner_id, source_intel_item_id, name, description, image_url, info
) VALUES (
    @base_id, @sector_x, @sector_y, @created_at, @type, @is_cloaked,
    @source_operation_id, @source_scanner_id, @source_intel_item_id, @name, @description, @image_url, @info
)
RETURNING id;

-- name: GetScanReportByID :one
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked,
       source_operation_id, source_scanner_id, source_intel_item_id, name, description, image_url, info
FROM game.scan_reports
WHERE id = @id;

-- name: RecentReportExistsByScanner :one
SELECT EXISTS (
    SELECT 1 FROM game.scan_reports
    WHERE source_scanner_id = @source_scanner_id
      AND created_at >= @since
);

-- name: ListScanReportsByBaseAndCoordinates :many
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked,
       source_operation_id, source_scanner_id, source_intel_item_id, name, description, image_url, info
FROM game.scan_reports
WHERE base_id = @base_id AND sector_x = @sector_x AND sector_y = @sector_y
ORDER BY created_at DESC;

-- name: GetLatestScanReportsByBase :many
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked,
       source_operation_id, source_scanner_id, source_intel_item_id, name, description, image_url, info
FROM game.scan_reports
WHERE base_id = @base_id
ORDER BY created_at DESC
LIMIT 50;

-- name: DeleteScanReport :exec
DELETE FROM game.scan_reports WHERE id = @id;

-- name: ListScanReportsByBaseWithinArea :many
WITH params AS (
    SELECT @base_id::bigint AS base_id,
           @center_x::int AS cx,
           @center_y::int AS cy,
           @radius::int AS r
)
SELECT sr.id, sr.base_id, sr.sector_x, sr.sector_y, sr.created_at, sr.type, sr.is_cloaked,
       sr.source_operation_id, sr.source_scanner_id, sr.source_intel_item_id, sr.name, sr.description, sr.image_url, sr.info
FROM game.scan_reports AS sr
JOIN game.sectors AS s ON s.x = sr.sector_x AND s.y = sr.sector_y
JOIN params p ON p.base_id = sr.base_id
WHERE ((s.x - p.cx) * (s.x - p.cx) + (s.y - p.cy) * (s.y - p.cy)) <= (p.r * p.r)
ORDER BY sr.created_at DESC;

