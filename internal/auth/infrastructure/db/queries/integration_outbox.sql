-- name: SaveIntegrationEvent :exec
INSERT INTO auth.integration_events (id, kind, payload, created_at, origin_id)
VALUES ($1, $2, $3, $4, $5);

-- name: IntegrationEventExists :one
SELECT EXISTS(SELECT 1 FROM auth.integration_events WHERE origin_id = $1 AND kind = $2);

-- name: ClaimUnpublishedIntegrationEvents :many
SELECT id, kind, payload, created_at
FROM auth.integration_events
WHERE published = FALSE
ORDER BY id ASC
FOR UPDATE SKIP LOCKED
LIMIT $1;

-- name: MarkIntegrationEventPublished :exec
UPDATE auth.integration_events
SET published = TRUE, published_at = $2
WHERE id = $1;
