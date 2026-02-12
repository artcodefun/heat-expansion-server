-- Sector queries

-- name: CreateSector :one
INSERT INTO game.sectors (x, y, name, description, image_url)
VALUES (@x, @y, @name, @description, @image_url)
RETURNING id, x, y, name, description, image_url;

-- name: UpdateSector :one
UPDATE game.sectors
SET name = @name,
    description = @description,
    image_url = @image_url
WHERE x = @x AND y = @y
RETURNING id, x, y, name, description, image_url;

-- name: GetSectorByCoordinates :one
SELECT id, x, y, name, description, image_url
FROM game.sectors
WHERE x = @x AND y = @y;

-- name: GetSectorByCoordinatesForUpdate :one
SELECT id, x, y, name, description, image_url
FROM game.sectors
WHERE x = @x AND y = @y
FOR UPDATE;

-- name: ListSectors :many
SELECT id, x, y, name, description, image_url
FROM game.sectors
ORDER BY id;

-- name: ListOccupiedSectorCoordinates :many
SELECT ub.sector_x AS x, ub.sector_y AS y FROM game.user_bases ub
UNION
SELECT rl.sector_x, rl.sector_y FROM game.resource_locations rl
UNION
SELECT dl.sector_x, dl.sector_y FROM game.dangerous_locations dl;

-- name: GetLocationTypeByCoordinates :one
SELECT CASE
    WHEN EXISTS (
        SELECT 1 FROM game.user_bases ub WHERE ub.sector_x = @x AND ub.sector_y = @y
    ) THEN 'BASE'
    WHEN EXISTS (
        SELECT 1 FROM game.resource_locations rl WHERE rl.sector_x = @x AND rl.sector_y = @y
    ) THEN 'RESOURCEFUL'
    WHEN EXISTS (
        SELECT 1 FROM game.dangerous_locations dl WHERE dl.sector_x = @x AND dl.sector_y = @y
    ) THEN 'DANGEROUS'
    ELSE 'EMPTY'
END AS location_type;

-- name: CountResourcefulLocationsInRange :one
SELECT COUNT(*) FROM game.resource_locations
WHERE sector_x >= (sqlc.arg('center_x')::int - sqlc.arg('radius')::int)
  AND sector_x <= (sqlc.arg('center_x')::int + sqlc.arg('radius')::int)
  AND sector_y >= (sqlc.arg('center_y')::int - sqlc.arg('radius')::int)
  AND sector_y <= (sqlc.arg('center_y')::int + sqlc.arg('radius')::int)
  AND (sector_x - sqlc.arg('center_x')::int)^2 + (sector_y - sqlc.arg('center_y')::int)^2 <= (sqlc.arg('radius')::int * sqlc.arg('radius')::int);

-- name: CountDangerousLocationsInRange :one
SELECT COUNT(*) FROM game.dangerous_locations
WHERE sector_x >= (sqlc.arg('center_x')::int - sqlc.arg('radius')::int)
  AND sector_x <= (sqlc.arg('center_x')::int + sqlc.arg('radius')::int)
  AND sector_y >= (sqlc.arg('center_y')::int - sqlc.arg('radius')::int)
  AND sector_y <= (sqlc.arg('center_y')::int + sqlc.arg('radius')::int)
  AND (sector_x - sqlc.arg('center_x')::int)^2 + (sector_y - sqlc.arg('center_y')::int)^2 <= (sqlc.arg('radius')::int * sqlc.arg('radius')::int);
