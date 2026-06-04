-- name: GetAdminProfile :one
SELECT id, username, active, created_at
FROM admin.admins
WHERE id = $1;
