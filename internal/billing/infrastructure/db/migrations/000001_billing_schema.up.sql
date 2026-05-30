CREATE SCHEMA billing;

CREATE TABLE billing.crystal_packages (
    id               UUID PRIMARY KEY,
    name             TEXT    NOT NULL,
    crystals         INT     NOT NULL,
    price_minor_units BIGINT NOT NULL,
    currency         TEXT    NOT NULL DEFAULT 'RUB',
    image_url        TEXT    NOT NULL DEFAULT '',
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       BIGINT  NOT NULL,
    updated_at       BIGINT  NOT NULL
);

CREATE TABLE billing.purchase_orders (
    id                UUID PRIMARY KEY,
    user_id           UUID    NOT NULL,
    package_id        UUID    NOT NULL REFERENCES billing.crystal_packages(id),
    crystals          INT     NOT NULL,
    amount_minor_units BIGINT NOT NULL,
    currency          TEXT    NOT NULL,
    provider          TEXT    NOT NULL,
    status            TEXT    NOT NULL DEFAULT 'PENDING',
    provider_order_id TEXT    NOT NULL DEFAULT '',
    confirmation_url  TEXT    NOT NULL DEFAULT '',
    created_at        BIGINT  NOT NULL,
    updated_at        BIGINT  NOT NULL
);

CREATE INDEX idx_billing_purchase_orders_user_id ON billing.purchase_orders (user_id);
CREATE INDEX idx_billing_purchase_orders_provider_order_id ON billing.purchase_orders (provider_order_id) WHERE provider_order_id <> '';

CREATE TABLE billing.domain_events (
    id           UUID    PRIMARY KEY,
    kind         TEXT    NOT NULL,
    payload      JSONB   NOT NULL,
    created_at   BIGINT  NOT NULL,
    published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at BIGINT
);

CREATE INDEX idx_billing_domain_events_published_id ON billing.domain_events (published, id);

CREATE TABLE billing.integration_events (
    id           UUID    PRIMARY KEY,
    kind         TEXT    NOT NULL,
    payload      JSONB   NOT NULL,
    created_at   BIGINT  NOT NULL,
    published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at BIGINT,
    origin_id    UUID,
    UNIQUE (origin_id, kind)
);

CREATE INDEX idx_billing_integration_events_published_id ON billing.integration_events (published, id);
