-- Outbox table for domain events (PostgreSQL)

-- Domain events outbox
CREATE TABLE game.domain_events (
    id           BIGSERIAL PRIMARY KEY,
    kind         TEXT   NOT NULL,
    payload      JSONB  NOT NULL,
    created_at   BIGINT NOT NULL,
    published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at BIGINT
);

CREATE INDEX idx_domain_events_published_id
    ON game.domain_events (published, id);
