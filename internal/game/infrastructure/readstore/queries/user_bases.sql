-- Base stats only (read repository no longer hydrates full overview)
-- name: GetBase :one
SELECT id, user_id, sector_x, sector_y, name, description, image_url
FROM game.user_bases
WHERE id = $1;

-- Base stats only (read repository no longer hydrates full overview)
-- name: GetBaseStats :one
SELECT stats, stats_calc_timestamp
FROM game.user_bases
WHERE id = $1;

-- List user-owned bases (basic info only)
-- name: ListUserBases :many
SELECT id, user_id, sector_x, sector_y, name, description, image_url
FROM game.user_bases
WHERE user_id = $1
ORDER BY id;

-- Base owner lookup by sector coordinates
-- name: GetBaseOwnerByCoordinates :one
SELECT u.id, u.name
FROM game.user_bases b
JOIN game.users u ON u.id = b.user_id
WHERE b.sector_x = $1 AND b.sector_y = $2;