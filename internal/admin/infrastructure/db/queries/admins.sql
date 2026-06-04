-- name: GetAdminByID :one
SELECT *
FROM admin.admins
WHERE id = $1;

-- name: GetAdminByUsername :one
SELECT *
FROM admin.admins
WHERE username = $1;

-- name: UpdateAdminCredentials :one
UPDATE admin.admins
SET password_hash     = @password_hash,
    invite_token      = @invite_token,
    updated_at        = @updated_at
WHERE id = @id
RETURNING *;
