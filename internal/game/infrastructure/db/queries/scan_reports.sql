-- Scan report queries

-- name: InsertScanReport :one
INSERT INTO game.scan_reports (
    base_id, sector_x, sector_y, created_at, type, is_cloaked,
    source_type, source_id,
    name, description, image_url, info
) VALUES (
    @base_id, @sector_x, @sector_y, @created_at, @type, @is_cloaked,
    @source_type, @source_id,
    @name, @description, @image_url, @info
)
RETURNING id;

-- name: GetScanReportByID :one
SELECT *
FROM game.scan_reports
WHERE id = @id;

-- name: RecentReportExistsByScanner :one
SELECT EXISTS (
    SELECT 1 FROM game.scan_reports
    WHERE source_type = 'SCANNER'
      AND source_id = @source_id
      AND created_at >= @since
);

-- name: ListScanReportsByBaseAndCoordinates :many
SELECT *
FROM game.scan_reports
WHERE base_id = @base_id AND sector_x = @sector_x AND sector_y = @sector_y
ORDER BY created_at DESC;

-- name: DeleteScanReport :exec
DELETE FROM game.scan_reports WHERE id = @id;

