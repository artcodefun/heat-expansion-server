-- Base stats only (read repository no longer hydrates full overview)
-- name: GetBaseStats :one
SELECT stats, stats_calc_timestamp
FROM user_bases
WHERE id = $1;