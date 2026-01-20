-- Prototypes catalog tables (static content referenced by aggregates)
-- Up migration: creates prototype tables and indexes.

-- Technologies catalog (must exist before others due to unlock_technology_id references)
CREATE TABLE tech_item_prototypes (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT    NOT NULL,
    category             TEXT    NOT NULL,
    unlock_technology_id BIGINT  REFERENCES tech_item_prototypes(id) ON DELETE RESTRICT,
    short_description    TEXT,
    full_description     TEXT,
    -- Price: {"credits": int, "iron": int, "titanium": int, "antimatter": int}
    price                JSONB   NOT NULL DEFAULT '{}'::jsonb,
    research_time        BIGINT  NOT NULL DEFAULT 0,
    image_url            TEXT,
    -- Improvement: {"type": string, "value": int, "max_level": int|null}
    improvement          JSONB
);
CREATE INDEX idx_tech_prototypes_category ON tech_item_prototypes(category);

-- Armies catalog
CREATE TABLE army_item_prototypes (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT    NOT NULL,
    category             TEXT    NOT NULL,
    faction              TEXT    NOT NULL DEFAULT 'EXO_COALITION',
    unlock_technology_id BIGINT  REFERENCES tech_item_prototypes(id) ON DELETE RESTRICT,
    short_description    TEXT,
    full_description     TEXT,
    -- Price: {"credits": int, "iron": int, "titanium": int, "antimatter": int}
    price                JSONB   NOT NULL DEFAULT '{}'::jsonb,
    production_time      BIGINT  NOT NULL DEFAULT 0,
    space                INTEGER NOT NULL DEFAULT 0,
    image_url            TEXT,
    attack               INTEGER NOT NULL DEFAULT 0,
    defence              INTEGER NOT NULL DEFAULT 0,
    capacity             INTEGER NOT NULL DEFAULT 0,
    stealth              INTEGER NOT NULL DEFAULT 0,
    speed                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_army_prototypes_category ON army_item_prototypes(category);

-- Buildings catalog (with category-specific JSONB payloads)
CREATE TABLE build_item_prototypes (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT    NOT NULL,
    category             TEXT    NOT NULL,
    faction              TEXT    NOT NULL DEFAULT 'EXO_COALITION',
    unlock_technology_id BIGINT  REFERENCES tech_item_prototypes(id) ON DELETE RESTRICT,
    short_description    TEXT,
    full_description     TEXT,
    -- Price: {"credits": int, "iron": int, "titanium": int, "antimatter": int}
    price                JSONB   NOT NULL DEFAULT '{}'::jsonb,
    production_time      BIGINT  NOT NULL DEFAULT 0,
    space                INTEGER NOT NULL DEFAULT 0,
    image_url            TEXT,
    -- Control building data: {"subtype": string}
    control_data         JSONB,
    -- Resources building data: {"credits_production": float, "iron_production": float, ..., "credits_capacity": int, ...}
    resources_data       JSONB,
    -- Defense building data: {"defence_bonus": int, "shield_strength": int}
    defense_data         JSONB,
    -- Military building data: {"unlock_army_category": string}
    military_data        JSONB,
    -- Intelligence building data: {"subtype": string, "stealth_strength": int, "target_location_type": string, "scan_range": int, "scan_cooldown": bigint}
    intelligence_data    JSONB
);
CREATE INDEX idx_build_prototypes_category ON build_item_prototypes(category);

-- Storage items catalog (with category-specific JSONB payloads)
CREATE TABLE storage_item_prototypes (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT    NOT NULL,
    category             TEXT    NOT NULL,
    estimated_worth      INTEGER NOT NULL DEFAULT 0,
    short_description    TEXT,
    full_description     TEXT,
    image_url            TEXT,
    -- Buff storage data: {"type": string, "value": float, "duration_seconds": bigint}
    buff_data            JSONB,
    -- Intel storage data: {"type": string, "decryption_seconds": bigint}
    intel_data           JSONB,
    -- Damaged storage data: {"restore_price": PriceDTO, "original_unit_id": int}
    damaged_data         JSONB,
    -- Artifact storage data: {"type": string, "value": float}
    artifact_data        JSONB,
    -- Consumable storage data: {"type": string, "box_contents": []string, "box_size": int}
    consumable_data      JSONB
);
CREATE INDEX idx_storage_prototypes_category ON storage_item_prototypes(category);
