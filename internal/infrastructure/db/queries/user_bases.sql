-- User bases queries

-- name: CreateBase :one
INSERT INTO user_bases (
    user_id, sector_x, sector_y, name, description, image_url,
    stats, stats_calc_timestamp
) VALUES (
    @user_id, @sector_x, @sector_y, @name, @description, @image_url,
    @stats, @stats_calc_timestamp
)
RETURNING id, user_id, sector_x, sector_y, name, description, image_url,
          stats, stats_calc_timestamp;

-- name: GetBaseByID :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE id = @id;

-- name: GetBaseByIDForUpdate :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE id = @id
FOR UPDATE;

-- name: UpdateBase :one
UPDATE user_bases
SET name = @name,
    description = @description,
    image_url = @image_url,
    stats = @stats,
    stats_calc_timestamp = @stats_calc_timestamp
WHERE id = @id
RETURNING id, user_id, sector_x, sector_y, name, description, image_url,
          stats, stats_calc_timestamp;
-- Note: above RETURNING still lists sector_id; fix below to sector_x, sector_y

-- name: DeleteBase :exec
DELETE FROM user_bases WHERE id = @id;

-- name: ListBasesByUserID :many
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE user_id = @user_id
ORDER BY id;

-- name: GetBaseByCoordinates :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE sector_x = @sector_x AND sector_y = @sector_y;

-- name: GetBaseByCoordinatesForUpdate :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE sector_x = @sector_x AND sector_y = @sector_y
FOR UPDATE;

-- name: FindClosestBase :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
WHERE sector_x != @x OR sector_y != @y
ORDER BY (sector_x - @x)^2 + (sector_y - @y)^2 ASC
LIMIT 1;

-- name: ListAllBases :many
SELECT id, user_id, sector_x, sector_y, name, description, image_url,
       stats, stats_calc_timestamp
FROM user_bases
ORDER BY id;
