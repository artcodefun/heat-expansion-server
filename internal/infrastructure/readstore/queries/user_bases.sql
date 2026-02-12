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