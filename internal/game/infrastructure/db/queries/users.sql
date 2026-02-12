-- Users queries

-- name: GetUserByID :one
SELECT id, name, email, password_hash, crystals
FROM game.users
WHERE id = @id;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, crystals
FROM game.users
WHERE email = @email;

-- name: InsertUser :one
INSERT INTO game.users (
    name, email, password_hash, crystals
) VALUES (
    @name, @email, @password_hash, @crystals
)
RETURNING id;

-- name: UpdateUser :exec
UPDATE game.users
SET name = @name,
    email = @email,
    password_hash = @password_hash,
    crystals = @crystals
WHERE id = @id;

-- name: ListUsers :many
SELECT id, name, email, password_hash, crystals
FROM game.users
ORDER BY id
LIMIT $1 OFFSET $2;
