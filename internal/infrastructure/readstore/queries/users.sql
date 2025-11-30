-- Readstore user profile queries

-- name: GetUserProfile :one
SELECT id, name, email, password_hash, crystals
FROM users
WHERE id = $1;
