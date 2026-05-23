CREATE TABLE game.black_market_offers (
    id                BIGSERIAL PRIMARY KEY,
    kind              TEXT    NOT NULL,
    prototype_id      BIGINT  NOT NULL,
    price_in_crystals BIGINT  NOT NULL,
    ends_at           BIGINT  NULL,
    is_limited        BOOLEAN NOT NULL DEFAULT FALSE,
    priority          BIGINT  NOT NULL DEFAULT 0,
    CONSTRAINT black_market_offers_kind_check CHECK (kind IN ('BUILDING', 'ARMY', 'STORAGE')),
    CONSTRAINT black_market_offers_unique_kind_prototype UNIQUE (kind, prototype_id)
);

CREATE INDEX idx_black_market_offers_kind ON game.black_market_offers(kind);
CREATE INDEX idx_black_market_offers_priority ON game.black_market_offers(priority DESC, id ASC);
CREATE INDEX idx_black_market_offers_limited_active ON game.black_market_offers(is_limited, ends_at);