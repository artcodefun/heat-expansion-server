-- Initial schema migration for Heat Expansion (PostgreSQL)
-- Up migration: creates all tables and indexes. No IF NOT EXISTS clauses.

-- Users
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT        NOT NULL,
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    crystals        INTEGER     NOT NULL DEFAULT 0
);

-- Sectors
CREATE TABLE sectors (
    id           BIGSERIAL PRIMARY KEY,
    x            INTEGER NOT NULL,
    y            INTEGER NOT NULL,
    name         TEXT,
    description  TEXT,
    image_url    TEXT,
    CONSTRAINT sectors_coordinates_unique UNIQUE (x, y)
);

-- User Bases (stats denormalized; internal items stored in separate tables)
CREATE TABLE user_bases (
    id                         BIGSERIAL PRIMARY KEY,
    user_id                    BIGINT    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sector_x                   INTEGER   NOT NULL,
    sector_y                   INTEGER   NOT NULL,
    name                       TEXT,
    description                TEXT,
    image_url                  TEXT,
    -- Base stats as JSONB (credits, iron, titanium, antimatter, capacities, production, attack/defence/space)
    stats                      JSONB            NOT NULL DEFAULT '{}'::jsonb,
    stats_calc_timestamp       BIGINT           NOT NULL DEFAULT 0,
    CONSTRAINT fk_user_bases_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES sectors(x, y) ON DELETE RESTRICT,
    CONSTRAINT user_bases_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_user_bases_user_id ON user_bases(user_id);

-- Military Operations
CREATE TABLE military_operations (
    id                 BIGSERIAL PRIMARY KEY,
    type               TEXT    NOT NULL,
    owner_user_id      BIGINT  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_base_id     BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
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
    -- Attacking units stored as JSONB array: [{"prototype_id": bigint, "count": int}]
    units              JSONB  NOT NULL DEFAULT '[]'::jsonb,
    spy_result         JSONB,
    attack_result      JSONB
);
CREATE INDEX idx_ops_owner ON military_operations(owner_user_id);
CREATE INDEX idx_ops_source_base ON military_operations(source_base_id);
CREATE INDEX idx_ops_target_coords ON military_operations(target_x, target_y);

-- Resource Locations (one per sector)
CREATE TABLE resource_locations (
    id               BIGSERIAL PRIMARY KEY,
    sector_x         INTEGER NOT NULL,
    sector_y         INTEGER NOT NULL,
    type             TEXT    NOT NULL,
    amount           INTEGER NOT NULL DEFAULT 0,
    name             TEXT,
    description      TEXT,
    image_url        TEXT,
    -- Resources as JSONB
    resources               JSONB  NOT NULL DEFAULT '{}'::jsonb,
    resources_calc_timestamp BIGINT  NOT NULL DEFAULT 0,
    -- Defenders as JSONB arrays: [{"prototype_id": bigint, "count": int}]
    units                   JSONB  NOT NULL DEFAULT '[]'::jsonb,
    structures              JSONB  NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT fk_resource_locations_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES sectors(x, y) ON DELETE CASCADE,
    CONSTRAINT resource_locations_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_resource_locations_sector_coords ON resource_locations(sector_x, sector_y);
-- Optional containment indexes for queries like units @> '[{"prototype_id": 123}]'
CREATE INDEX idx_resource_locations_units_gin ON resource_locations USING gin (units jsonb_path_ops);
CREATE INDEX idx_resource_locations_structures_gin ON resource_locations USING gin (structures jsonb_path_ops);

-- Dangerous Locations (one per sector)
CREATE TABLE dangerous_locations (
    id               BIGSERIAL PRIMARY KEY,
    sector_x         INTEGER NOT NULL,
    sector_y         INTEGER NOT NULL,
    danger_level     INTEGER NOT NULL DEFAULT 0,
    name             TEXT,
    description      TEXT,
    image_url        TEXT,
    -- Resources as JSONB
    resources               JSONB  NOT NULL DEFAULT '{}'::jsonb,
    resources_calc_timestamp BIGINT  NOT NULL DEFAULT 0,
    -- Defenders as JSONB arrays: [{"prototype_id": bigint, "count": int}]
    units                   JSONB  NOT NULL DEFAULT '[]'::jsonb,
    structures              JSONB  NOT NULL DEFAULT '[]'::jsonb,
    CONSTRAINT fk_dangerous_locations_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES sectors(x, y) ON DELETE CASCADE,
    CONSTRAINT dangerous_locations_sector_unique UNIQUE (sector_x, sector_y)
);
CREATE INDEX idx_dangerous_locations_sector_coords ON dangerous_locations(sector_x, sector_y);
-- Optional containment indexes
CREATE INDEX idx_dang_locations_units_gin ON dangerous_locations USING gin (units jsonb_path_ops);
CREATE INDEX idx_dang_locations_structures_gin ON dangerous_locations USING gin (structures jsonb_path_ops);

-- Empty Locations table removed; emptiness is derived.

-- Sector Scan Reports
CREATE TABLE scan_reports (
    id                   BIGSERIAL PRIMARY KEY,
    base_id              BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    sector_x             INTEGER NOT NULL,
    sector_y             INTEGER NOT NULL,
    created_at           BIGINT  NOT NULL,
    type                 TEXT    NOT NULL,
    is_cloaked           BOOLEAN NOT NULL DEFAULT FALSE,
    source_operation_id  BIGINT  REFERENCES military_operations(id) ON DELETE SET NULL,
    name                 TEXT,
    description          TEXT,
    image_url            TEXT,
    -- Info as JSONB
    info                JSONB  NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT fk_scan_reports_sector_coords FOREIGN KEY (sector_x, sector_y) REFERENCES sectors(x, y) ON DELETE CASCADE
);
CREATE INDEX idx_scan_reports_base_created_at ON scan_reports(base_id, created_at DESC);
CREATE INDEX idx_scan_reports_sector_coords ON scan_reports(sector_x, sector_y);

-- Activities (append-only feed; payload captures subtype-specific data)
CREATE TABLE activities (
    id          BIGSERIAL PRIMARY KEY,
    kind        TEXT   NOT NULL,
    created_at  BIGINT NOT NULL,
    base_id     BIGINT NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    -- Per-kind payloads as JSONB (only one is non-NULL according to kind)
    operation_data   JSONB,
    scan_data        JSONB,
    radar_data       JSONB,
    trade_data       JSONB,
    CONSTRAINT chk_activity_kind CHECK (kind IN ('MILITARY','SCAN','RADAR','TRADE')),
    CONSTRAINT chk_activity_payload_by_kind CHECK (
        (kind = 'MILITARY' AND operation_data IS NOT NULL AND scan_data IS NULL AND radar_data IS NULL AND trade_data IS NULL) OR
        (kind = 'SCAN'     AND operation_data IS NULL AND scan_data IS NOT NULL AND radar_data IS NULL AND trade_data IS NULL) OR
        (kind = 'RADAR'    AND operation_data IS NULL AND scan_data IS NULL AND radar_data IS NOT NULL AND trade_data IS NULL) OR
        (kind = 'TRADE'    AND operation_data IS NULL AND scan_data IS NULL AND radar_data IS NULL AND trade_data IS NOT NULL)
    )
);
CREATE INDEX idx_activities_base_created_at ON activities(base_id, created_at DESC);
