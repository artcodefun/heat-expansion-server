-- name: CreateAccount :one
INSERT INTO auth.users (id, name, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, password_hash;

-- name: GetAccountByEmail :one
SELECT id, name, email, password_hash
FROM auth.users
WHERE email = $1;

-- name: GetAccountByID :one
SELECT id, name, email, password_hash
FROM auth.users
WHERE id = $1;

-- name: UpdateAccountPasswordHash :exec
UPDATE auth.users
SET password_hash = $2
WHERE id = $1;
