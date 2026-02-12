-- Outbox domain events queries

-- name: InsertOutboxEvent :one
INSERT INTO game.domain_events (
    kind, payload, created_at, published
) VALUES (
    @kind, @payload, @created_at, FALSE
)
RETURNING id;

-- name: ClaimUnpublishedOutboxEvents :many
SELECT id, kind, payload, created_at, published, published_at
FROM game.domain_events
WHERE published = FALSE
ORDER BY id ASC
FOR UPDATE SKIP LOCKED
LIMIT $1;

-- name: MarkOutboxEventPublished :exec
UPDATE game.domain_events
SET published = TRUE,
    published_at = @published_at
WHERE id = @id;
