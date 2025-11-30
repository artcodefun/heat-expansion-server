-- Sector and scan reports queries

-- name: GetSector :one
SELECT x, y, name, description, image_url
FROM sectors
WHERE x = $1 AND y = $2;

-- name: GetLatestScans :many
SELECT DISTINCT ON (sector_x, sector_y) id, base_id, sector_x, sector_y, created_at, type, is_cloaked, source_operation_id,
       name, description, image_url, info
FROM scan_reports
WHERE base_id = $1
ORDER BY sector_x, sector_y, created_at DESC;

-- name: GetScansNear :many
SELECT id, base_id, sector_x, sector_y, created_at, type, is_cloaked, source_operation_id,
       name, description, image_url, info
FROM scan_reports
WHERE base_id = $1 AND ((sector_x - $2) * (sector_x - $2) + (sector_y - $3) * (sector_y - $3)) <= ($4 * $4)
ORDER BY created_at DESC;

-- name: ListOccupiedCoordinates :many
SELECT sector_x AS x, sector_y AS y FROM user_bases
UNION
SELECT sector_x AS x, sector_y AS y FROM resource_locations
UNION
SELECT sector_x AS x, sector_y AS y FROM dangerous_locations;

-- name: ListSectorsInRadius :many
SELECT x, y, name, description, image_url
FROM sectors
WHERE ((x - $1) * (x - $1) + (y - $2) * (y - $2)) <= ($3 * $3);
