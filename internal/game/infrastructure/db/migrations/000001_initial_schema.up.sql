-- Initial schema migration for Heat Expansion (PostgreSQL)
-- Up migration: creates all tables and indexes. No IF NOT EXISTS clauses.

CREATE SCHEMA game;

-- Users
CREATE TABLE game.users (
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT        NOT NULL,
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    crystals        INTEGER     NOT NULL DEFAULT 0
);

-- Sectors
CREATE TABLE game.sectors (
    id           BIGSERIAL PRIMARY KEY,
    x            INTEGER NOT NULL,
    y            INTEGER NOT NULL,
    name         TEXT,
    description  TEXT,
    image_url    TEXT,
    CONSTRAINT sectors_coordinates_unique UNIQUE (x, y)
);

-- User Bases (stats denormalized; internal items stored in separate tables)
CREATE TABLE game.user_bases (
    id                         BIGSERIAL PRIMARY KEY,
    user_id                    BIGINT    NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    sector_x                   INTEGER   NOT NULL,
    sector_y                   INTEGER   NOT NULL,
    name                       TEXT,
    description                TEXT,
    image_url                  TEXT,
    -- Base stats: {"credits": int, "iron": int, "titanium": int, "antimatter": int, "defence": int, "attack": int, "space": int, ...}
    stats                      JSONB            NOT NULL DEFAULT '{}'::jsonb,
    stats_calc_timestamp       BIGINT           NOT NULL DEFAULT 0,
    CONSTRAINT fk_user_bases_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES game.sectors(x, y) ON DELETE RESTRICT,
    CONSTRAINT user_bases_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_user_bases_user_id ON game.user_bases(user_id);

-- Military Operations
CREATE TABLE game.military_operations (
    id                 BIGSERIAL PRIMARY KEY,
    type               TEXT    NOT NULL,
    owner_user_id      BIGINT  NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    source_base_id     BIGINT  NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    -- target sector is identified by coordinates; keep id-less reference
    source_x           INTEGER NOT NULL,
    source_y           INTEGER NOT NULL,
    target_x           INTEGER NOT NULL,
    target_y           INTEGER NOT NULL,
    outbound_depart_at BIGINT  NOT NULL DEFAULT 0,
    outbound_arrive_at BIGINT  NOT NULL DEFAULT 0,
    return_depart_at   BIGINT  NOT NULL DEFAULT 0,
    return_arrive_at   BIGINT  NOT NULL DEFAULT 0,
    completed_at       BIGINT  NOT NULL DEFAULT 0,
    phase              TEXT    NOT NULL,
    result             TEXT    NOT NULL,
    crystals_skip_price INTEGER NOT NULL DEFAULT 0,
    -- Snapshot of units: [{"prototype_id": int, "category": string, "attack": int, "defence": int, "count": int, ...}]
    units              JSONB  NOT NULL DEFAULT '[]'::jsonb,
    -- Snapshot of active buffs/artifacts at operation creation: [{"prototype_id": int, "category": string, "buff_data": {...}, "artifact_data": {...}}]
    storage_snaps      JSONB  NOT NULL DEFAULT '[]'::jsonb,
    -- Pre-calculated total multipliers for the attacker: {"attack_mul": float, "defence_mul": float, "stealth_mul": float, "capacity_mul": float, "speed_mul": float}
    total_modifiers    JSONB  NOT NULL DEFAULT '{}'::jsonb,
    -- Spy result: {"outcome": string, "attacker_remaining": [...], "defender_remaining": [...], "defenders_before": [...]}
    spy_result         JSONB,
    -- Attack result: {"outcome": string, "attacker_remaining": [...], "defender_remaining": [...], "remaining_structures": [...], "loot": {...}, ...}
    attack_result      JSONB
);
CREATE INDEX idx_ops_owner ON game.military_operations(owner_user_id);
CREATE INDEX idx_ops_source_base ON game.military_operations(source_base_id);
CREATE INDEX idx_ops_target_coords ON game.military_operations(target_x, target_y);

-- Resource Locations (one per sector)
CREATE TABLE game.resource_locations (
    id                BIGSERIAL PRIMARY KEY,
    sector_x          INTEGER NOT NULL,
    sector_y          INTEGER NOT NULL,
    resource_type     TEXT    NOT NULL,
    defender_faction  TEXT    NOT NULL,
    total_worth       INTEGER NOT NULL DEFAULT 0,
    name              TEXT,
    description       TEXT,
    image_url         TEXT,
    -- Resources: {"credits": int, "iron": int, "titanium": int, "antimatter": int}
    resources               JSONB  NOT NULL DEFAULT '{}'::jsonb,
    resources_calc_timestamp BIGINT  NOT NULL DEFAULT 0,
    -- Defenders: [{"prototype_id": int, "count": int}]
    armies                   JSONB  NOT NULL DEFAULT '[]'::jsonb,
    -- Buildings: [{"prototype_id": int, "count": int}]
    buildings                JSONB  NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT fk_resource_locations_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES game.sectors(x, y) ON DELETE CASCADE,
    CONSTRAINT resource_locations_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_resource_locations_sector_coords ON game.resource_locations(sector_x, sector_y);
-- Optional containment indexes for queries like armies @> '[{"prototype_id": 123}]'
CREATE INDEX idx_resource_locations_armies_gin ON game.resource_locations USING gin (armies jsonb_path_ops);
CREATE INDEX idx_resource_locations_buildings_gin ON game.resource_locations USING gin (buildings jsonb_path_ops);

-- Dangerous Locations (one per sector)
CREATE TABLE game.dangerous_locations (
    id                BIGSERIAL PRIMARY KEY,
    sector_x          INTEGER NOT NULL,
    sector_y          INTEGER NOT NULL,
    defender_faction  TEXT    NOT NULL,
    total_worth       INTEGER NOT NULL DEFAULT 0,
    name              TEXT,
    description       TEXT,
    image_url         TEXT,
    -- Resources: {"credits": int, "iron": int, "titanium": int, "antimatter": int}
    resources               JSONB  NOT NULL DEFAULT '{}'::jsonb,
    resources_calc_timestamp BIGINT  NOT NULL DEFAULT 0,
    -- Defenders: [{"prototype_id": int, "count": int}]
    armies                   JSONB  NOT NULL DEFAULT '[]'::jsonb,
    -- Buildings: [{"prototype_id": int, "count": int}]
    buildings                JSONB  NOT NULL DEFAULT '[]'::jsonb,
    -- Trophies: [{"prototype_id": int}]
    trophies                 JSONB  NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT fk_dangerous_locations_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES game.sectors(x, y) ON DELETE CASCADE,
    CONSTRAINT dangerous_locations_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_dangerous_locations_sector_coords ON game.dangerous_locations(sector_x, sector_y);
-- Optional containment indexes
CREATE INDEX idx_dang_locations_armies_gin ON game.dangerous_locations USING gin (armies jsonb_path_ops);
CREATE INDEX idx_dang_locations_buildings_gin ON game.dangerous_locations USING gin (buildings jsonb_path_ops);

-- Sector Scan Reports
CREATE TABLE game.scan_reports (
    id                   BIGSERIAL PRIMARY KEY,
    base_id              BIGINT  NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    sector_x             INTEGER NOT NULL,
    sector_y             INTEGER NOT NULL,
    created_at           BIGINT  NOT NULL,
    type                 TEXT    NOT NULL,
    is_cloaked           BOOLEAN NOT NULL DEFAULT FALSE,
    source_operation_id  BIGINT  REFERENCES game.military_operations(id) ON DELETE SET NULL,
    source_scanner_id    UUID,
    source_intel_item_id UUID,
    name                 TEXT,
    description          TEXT,
    image_url            TEXT,
    -- Scan info: {"credits": int, "iron": int, "titanium": int, "antimatter": int, "defence": int, "attack": int, "space": int}
    info                JSONB  NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT fk_scan_reports_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES game.sectors(x, y) ON DELETE CASCADE
);
CREATE INDEX idx_scan_reports_base_created_at ON game.scan_reports(base_id, created_at DESC);
CREATE INDEX idx_scan_reports_sector_coords ON game.scan_reports(sector_x, sector_y);

-- Activities (append-only feed; payload captures subtype-specific data)
CREATE TABLE game.activities (
    id          UUID PRIMARY KEY,
    kind        TEXT   NOT NULL,
    created_at  BIGINT NOT NULL,
    base_id     BIGINT NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    -- Offense activity: {"op_id": int, "subtype": string}
    offense_data     JSONB,
    -- Defense activity: {"op_id": int, "subtype": string}
    defense_data     JSONB,
    -- Scan activity: {"report_id": int}
    scan_data        JSONB,
    -- Radar activity: {"op_id": int, "detected_at": bigint, "eta_at_base": bigint, "source_x": int, "source_y": int, "target_x": int, "target_y": int, "threat": {"attack": int, "speed": int, ...}}
    radar_data       JSONB,
    -- Trade activity: JSONB structure TBD
    trade_data       JSONB,
    CONSTRAINT chk_activity_kind CHECK (kind IN ('OFFENSE','DEFENSE','SCAN','RADAR','TRADE')),
    CONSTRAINT chk_activity_payload_by_kind CHECK (
        (kind = 'OFFENSE' AND offense_data IS NOT NULL AND defense_data IS NULL AND scan_data IS NULL AND radar_data IS NULL AND trade_data IS NULL) OR
        (kind = 'DEFENSE' AND offense_data IS NULL AND defense_data IS NOT NULL AND scan_data IS NULL AND radar_data IS NULL AND trade_data IS NULL) OR
        (kind = 'SCAN'     AND offense_data IS NULL AND defense_data IS NULL AND scan_data IS NOT NULL AND radar_data IS NULL AND trade_data IS NULL) OR
        (kind = 'RADAR'    AND offense_data IS NULL AND defense_data IS NULL AND scan_data IS NULL AND radar_data IS NOT NULL AND trade_data IS NULL) OR
        (kind = 'TRADE'    AND offense_data IS NULL AND defense_data IS NULL AND scan_data IS NULL AND radar_data IS NULL AND trade_data IS NOT NULL)
    )
);
CREATE INDEX idx_activities_base_created_at ON game.activities(base_id, created_at DESC);
