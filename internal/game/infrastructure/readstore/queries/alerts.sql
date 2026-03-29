-- name: ListAlertsByUser :many
SELECT * FROM game.alerts
WHERE user_id = $1 AND expires_at > $2
ORDER BY created_at DESC;

-- name: CountUnreadAlertsByUser :one
SELECT count(*) FROM game.alerts
WHERE user_id = $1 AND is_read = false AND expires_at > $2;
