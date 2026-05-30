-- name: SaveOutboxEvent :exec
INSERT INTO billing.domain_events (id, kind, payload, created_at)
VALUES ($1, $2, $3, $4);

-- name: ClaimUnpublishedEvents :many
SELECT id, kind, payload, created_at
FROM billing.domain_events
WHERE published = FALSE
ORDER BY id ASC
FOR UPDATE SKIP LOCKED
LIMIT $1;

-- name: MarkEventPublished :exec
UPDATE billing.domain_events
SET published = TRUE, published_at = $2
WHERE id = $1;

-- name: NotifyOutboxEvent :exec
NOTIFY billing_domain_events;
