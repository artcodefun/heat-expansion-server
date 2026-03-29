-- name: InsertAlert :exec
INSERT INTO game.alerts (id, user_id, base_id, activity_id, kind, title, content, is_read, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: ListAlertsByUser :many
SELECT * FROM game.alerts
WHERE user_id = $1 AND expires_at > $2
ORDER BY created_at DESC;

-- name: DeleteExpiredAlerts :exec
DELETE FROM game.alerts WHERE expires_at < $1;

-- name: MarkAllAlertsAsReadByUser :exec
UPDATE game.alerts SET is_read = TRUE
WHERE user_id = $1;

-- name: ExistsForActivity :one
SELECT EXISTS (
    SELECT 1 FROM game.alerts
    WHERE activity_id = $1
);
