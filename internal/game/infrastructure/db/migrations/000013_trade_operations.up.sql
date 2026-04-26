CREATE TABLE game.trade_operations (
    id                    BIGSERIAL PRIMARY KEY,
    operation_uuid        UUID    NOT NULL,
    created_at            BIGINT  NOT NULL DEFAULT 0,
    sender_user_id        UUID    NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    sender_base_id        BIGINT  NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    receiver_user_id      UUID    NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    receiver_base_id      BIGINT  NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    source_x              INTEGER NOT NULL,
    source_y              INTEGER NOT NULL,
    target_x              INTEGER NOT NULL,
    target_y              INTEGER NOT NULL,
    offered_payload       JSONB   NOT NULL DEFAULT '{}'::jsonb,
    requested_payload     JSONB   NOT NULL DEFAULT '{}'::jsonb,
    transport_units       JSONB   NOT NULL DEFAULT '[]'::jsonb,
    storage_snaps         JSONB   NOT NULL DEFAULT '[]'::jsonb,
    total_modifiers       JSONB   NOT NULL DEFAULT '{}'::jsonb,
    expires_at            BIGINT  NOT NULL,
    outbound_depart_at    BIGINT  NOT NULL DEFAULT 0,
    outbound_arrive_at    BIGINT  NOT NULL DEFAULT 0,
    arrived_at_target_at  BIGINT  NOT NULL DEFAULT 0,
    return_depart_at      BIGINT  NOT NULL DEFAULT 0,
    return_arrive_at      BIGINT  NOT NULL DEFAULT 0,
    completed_at          BIGINT  NOT NULL DEFAULT 0,
    phase                 TEXT    NOT NULL,
    result                TEXT    NOT NULL,
    crystals_skip_price   BIGINT  NOT NULL DEFAULT 0,
    CONSTRAINT trade_operations_uuid_unique UNIQUE (operation_uuid)
);

CREATE INDEX idx_trade_ops_sender_base ON game.trade_operations(sender_base_id);
CREATE INDEX idx_trade_ops_receiver_base ON game.trade_operations(receiver_base_id);
CREATE INDEX idx_trade_ops_phase ON game.trade_operations(phase);
CREATE INDEX idx_trade_ops_expires_at ON game.trade_operations(expires_at);