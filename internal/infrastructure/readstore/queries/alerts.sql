-- name: ListAlertsByBase :many
SELECT * FROM alerts
WHERE base_id = $1 AND expires_at > $2
ORDER BY created_at DESC;

-- name: CountUnreadAlertsByBase :one
SELECT count(*) FROM alerts
WHERE base_id = $1 AND is_read = false AND expires_at > $2;
