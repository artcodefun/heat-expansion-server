-- Readstore user profile queries

-- name: GetUserProfile :one
SELECT id, name, crystals
FROM game.users
WHERE id = $1;
