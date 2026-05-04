-- name: CreatePasswordResetToken :exec
INSERT INTO auth.password_reset_tokens (id, account_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetActivePasswordResetToken :one
SELECT id, account_id, token_hash, expires_at, used_at
FROM auth.password_reset_tokens
WHERE account_id = $1 AND token_hash = $2 AND used_at IS NULL AND expires_at > EXTRACT(EPOCH FROM NOW())::BIGINT;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE auth.password_reset_tokens
SET used_at = $2
WHERE id = $1;

-- name: InvalidateAccountPasswordResetTokens :exec
UPDATE auth.password_reset_tokens
SET used_at = $2
WHERE account_id = $1 AND used_at IS NULL;
