-- Users queries

-- name: GetUserByID :one
SELECT id, name, crystals
FROM game.users
WHERE id = @id;

-- name: GetUserByIDForUpdate :one
SELECT id, name, crystals
FROM game.users
WHERE id = @id
FOR UPDATE;

-- name: InsertUser :exec
INSERT INTO game.users (
    id, name, crystals
) VALUES (
    @id, @name, @crystals
);

-- name: UpdateUser :exec
UPDATE game.users
SET name = @name,
    crystals = @crystals
WHERE id = @id;

-- name: ListUsers :many
SELECT id, name, crystals
FROM game.users
ORDER BY id
LIMIT $1 OFFSET $2;
