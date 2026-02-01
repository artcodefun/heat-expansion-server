-- name: InsertAlert :exec
INSERT INTO alerts (id, base_id, activity_id, kind, title, content, is_read, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: ListAlertsByBase :many
SELECT * FROM alerts
WHERE base_id = $1 AND expires_at > $2
ORDER BY created_at DESC;

-- name: DeleteExpiredAlerts :exec
DELETE FROM alerts WHERE expires_at < $1;

-- name: MarkAllAlertsAsRead :exec
UPDATE alerts SET is_read = TRUE
WHERE base_id = $1;

-- name: ExistsForActivity :one
SELECT EXISTS (
    SELECT 1 FROM alerts
    WHERE activity_id = $1
);
