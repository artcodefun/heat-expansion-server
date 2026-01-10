-- Dangerous locations queries

-- name: GetDangerousLocationByID :one
SELECT id, sector_x, sector_y, danger_level, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings, trophies
FROM dangerous_locations
WHERE id = @id;

-- name: GetDangerousLocationBySector :one
SELECT id, sector_x, sector_y, danger_level, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings, trophies
FROM dangerous_locations
WHERE sector_x = @sector_x AND sector_y = @sector_y;

-- name: GetDangerousLocationBySectorForUpdate :one
SELECT id, sector_x, sector_y, danger_level, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings, trophies
FROM dangerous_locations
WHERE sector_x = @sector_x AND sector_y = @sector_y
FOR UPDATE;

-- name: FindClosestDangerousLocation :one
SELECT id, sector_x, sector_y, danger_level, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings, trophies
FROM dangerous_locations
ORDER BY (sector_x - @x)^2 + (sector_y - @y)^2 ASC
LIMIT 1;

-- name: InsertDangerousLocation :one
INSERT INTO dangerous_locations (
    sector_x, sector_y, danger_level, name, description, image_url,
    resources, resources_calc_timestamp, armies, buildings, trophies
) VALUES (
    @sector_x, @sector_y, @danger_level, @name, @description, @image_url,
    @resources, @resources_calc_timestamp, @armies, @buildings, @trophies
)
RETURNING id;

-- name: UpdateDangerousLocation :exec
UPDATE dangerous_locations
SET danger_level = @danger_level,
    name = @name,
    description = @description,
    image_url = @image_url,
    resources = @resources,
    resources_calc_timestamp = @resources_calc_timestamp,
    armies = @armies,
    buildings = @buildings,
    trophies = @trophies
WHERE id = @id;

-- name: DeleteDangerousLocation :exec
DELETE FROM dangerous_locations WHERE id = @id;
