-- Radar Threats table for live tracking of incoming hostilities.
CREATE TABLE radar_threats (
    id                 UUID          PRIMARY KEY,
    operation_id       BIGINT        NOT NULL REFERENCES military_operations(id) ON DELETE CASCADE,
    owner_base_id      BIGINT        NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    detected_at        BIGINT        NOT NULL,
    source_x           INTEGER       NOT NULL,
    source_y           INTEGER       NOT NULL,
    target_x           INTEGER       NOT NULL,
    target_y           INTEGER       NOT NULL,
    estimated_arrival_at BIGINT      NOT NULL,
    arrival_at         BIGINT,
    type               TEXT          NOT NULL,
    status             TEXT          NOT NULL,
    attack             INTEGER       NOT NULL,
    speed              INTEGER       NOT NULL,
    stealth            INTEGER       NOT NULL,
    capacity           INTEGER       NOT NULL,
    CONSTRAINT chk_radar_threat_status CHECK (status IN ('ARRIVING', 'LOST', 'ARRIVED')),
    CONSTRAINT chk_radar_threat_type CHECK (type IN ('ATTACK', 'SPY'))
);

CREATE INDEX idx_radar_threats_owner_base_status ON radar_threats(owner_base_id, status);
CREATE INDEX idx_radar_threats_operation_id ON radar_threats(operation_id);
