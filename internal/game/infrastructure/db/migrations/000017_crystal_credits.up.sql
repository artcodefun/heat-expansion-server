CREATE TABLE IF NOT EXISTS game.crystal_credits (
    order_id    UUID        PRIMARY KEY,
    user_id     UUID        NOT NULL,
    crystals    INT         NOT NULL,
    credited_at BIGINT      NOT NULL
);
