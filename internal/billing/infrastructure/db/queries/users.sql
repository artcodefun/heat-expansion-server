-- name: UpsertUser :exec
INSERT INTO billing.users (id, email)
VALUES (@id, @email)
ON CONFLICT (id) DO UPDATE SET email = EXCLUDED.email;

-- name: GetUserByID :one
SELECT id, email
FROM billing.users
WHERE id = @id;
