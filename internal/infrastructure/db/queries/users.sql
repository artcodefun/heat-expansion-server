-- Users queries

-- name: GetUserByID :one
SELECT id, name, email, password_hash, crystals
FROM users
WHERE id = @id;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, crystals
FROM users
WHERE email = @email;

-- name: InsertUser :one
INSERT INTO users (
    name, email, password_hash, crystals
) VALUES (
    @name, @email, @password_hash, @crystals
)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET name = @name,
    email = @email,
    password_hash = @password_hash,
    crystals = @crystals
WHERE id = @id;

-- name: ListUsers :many
SELECT id, name, email, password_hash, crystals
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;
