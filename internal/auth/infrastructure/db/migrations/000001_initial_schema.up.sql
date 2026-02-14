-- Initial schema for Auth service

CREATE SCHEMA auth;

CREATE TABLE auth.users (
    id              UUID PRIMARY KEY,
    name            TEXT        NOT NULL,
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- Outbox table for auth domain events
CREATE TABLE auth.domain_events (
    id           UUID PRIMARY KEY,
    kind         TEXT   NOT NULL,
    payload      JSONB  NOT NULL,
    created_at   BIGINT NOT NULL,
    published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at BIGINT
);

CREATE INDEX idx_auth_domain_events_published_id
    ON auth.domain_events (published, id);

-- Integration outbox table for auth service
CREATE TABLE auth.integration_events (
    id           UUID PRIMARY KEY,
    kind         TEXT    NOT NULL,
    payload      JSONB   NOT NULL,
    created_at   BIGINT  NOT NULL,
    published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at BIGINT,
    origin_id    UUID,
    UNIQUE (origin_id, kind)
);

CREATE INDEX idx_auth_integration_events_published_id
    ON auth.integration_events (published, id);
