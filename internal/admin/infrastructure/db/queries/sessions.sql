-- name: CreateSession :exec
INSERT INTO admin.sessions (token, admin_id, expires_at, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetSessionByToken :one
SELECT *
FROM admin.sessions
WHERE token = @token
  AND expires_at > @now;

-- name: DeleteSession :exec
DELETE FROM admin.sessions
WHERE token = $1;
