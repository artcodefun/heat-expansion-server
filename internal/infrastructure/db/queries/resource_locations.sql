-- Resource locations queries

-- name: GetResourceLocationByID :one
SELECT id, sector_x, sector_y, type, amount, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings
FROM resource_locations
WHERE id = @id;

-- name: GetResourceLocationBySector :one
SELECT id, sector_x, sector_y, type, amount, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings
FROM resource_locations
WHERE sector_x = @sector_x AND sector_y = @sector_y;

-- name: GetResourceLocationBySectorForUpdate :one
SELECT id, sector_x, sector_y, type, amount, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings
FROM resource_locations
WHERE sector_x = @sector_x AND sector_y = @sector_y
FOR UPDATE;

-- name: ListResourceLocations :many
SELECT id, sector_x, sector_y, type, amount, name, description, image_url,
       resources, resources_calc_timestamp, armies, buildings
FROM resource_locations
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: InsertResourceLocation :one
INSERT INTO resource_locations (
    sector_x, sector_y, type, amount, name, description, image_url,
    resources, resources_calc_timestamp, armies, buildings
) VALUES (
    @sector_x, @sector_y, @type, @amount, @name, @description, @image_url,
    @resources, @resources_calc_timestamp, @armies, @buildings
)
RETURNING id;

-- name: UpdateResourceLocation :exec
UPDATE resource_locations
SET type = @type,
    amount = @amount,
    name = @name,
    description = @description,
    image_url = @image_url,
    resources = @resources,
    resources_calc_timestamp = @resources_calc_timestamp,
    armies = @armies,
    buildings = @buildings
WHERE id = @id;

-- name: DeleteResourceLocation :exec
DELETE FROM resource_locations WHERE id = @id;
