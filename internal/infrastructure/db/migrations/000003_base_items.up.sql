-- Base items tables for normalized item storage per base
-- Up migration: creates per-item tables with status columns, constraints, and indexes.

-- Ensure pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE base_army_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id             BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    prototype_id        BIGINT  NOT NULL REFERENCES army_item_prototypes(id) ON DELETE RESTRICT,
    status              TEXT    NOT NULL,
    -- Army pending data: {"count": int}
    pending_data        JSONB,
    -- Army in production data: {"start_date": bigint, "completion_date": bigint, "crystals_skip_price": int}
    in_prod_data        JSONB,
    -- Army present data: {"count": int, "refund": {...}}
    present_data        JSONB,
    -- Army deployed data: {"operation_id": int, "count": int}
    deployed_data       JSONB,
    -- bookkeeping
    created_at          BIGINT  NOT NULL,
    -- constraints
    CONSTRAINT chk_army_status CHECK (status IN ('PENDING','IN_PRODUCTION','PRESENT','DEPLOYED')),
    CONSTRAINT chk_army_pending_json CHECK (
        status <> 'PENDING' OR pending_data IS NOT NULL
    ),
    CONSTRAINT chk_army_in_production_json CHECK (
        status <> 'IN_PRODUCTION' OR in_prod_data IS NOT NULL
    ),
    CONSTRAINT chk_army_present_json CHECK (
        status <> 'PRESENT' OR present_data IS NOT NULL
    ),
    CONSTRAINT chk_army_deployed_json CHECK (
        status <> 'DEPLOYED' OR deployed_data IS NOT NULL
    )
);
CREATE INDEX idx_base_army_items_base ON base_army_items(base_id);
CREATE INDEX idx_base_army_items_base_status ON base_army_items(base_id, status);
CREATE INDEX idx_base_army_items_due ON base_army_items(((in_prod_data->>'completion_date')::bigint)) WHERE status = 'IN_PRODUCTION';
CREATE UNIQUE INDEX uq_base_army_items_present ON base_army_items(base_id, prototype_id) WHERE status = 'PRESENT';
CREATE UNIQUE INDEX uq_base_army_items_pending ON base_army_items(base_id, prototype_id) WHERE status = 'PENDING';
CREATE INDEX idx_base_army_items_operation ON base_army_items(((deployed_data->>'operation_id')::bigint)) WHERE status = 'DEPLOYED';

CREATE TABLE base_build_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id             BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    prototype_id        BIGINT  NOT NULL REFERENCES build_item_prototypes(id) ON DELETE RESTRICT,
    status              TEXT    NOT NULL,
    -- Build pending data: {}
    pending_data        JSONB,
    -- Build in production data: {"start_date": bigint, "completion_date": bigint, "crystals_skip_price": int}
    in_prod_data        JSONB,
    -- Build present data: {"refund": {...}}
    present_data        JSONB,
    -- bookkeeping
    created_at          BIGINT  NOT NULL,
    -- constraints
    CONSTRAINT chk_build_status CHECK (status IN ('PENDING','IN_PRODUCTION','PRESENT')),
    CONSTRAINT chk_build_in_production_json CHECK (
        status <> 'IN_PRODUCTION' OR in_prod_data IS NOT NULL
    ),
    CONSTRAINT chk_build_pending_json CHECK (
        status <> 'PENDING' OR pending_data IS NOT NULL
    )
);
CREATE INDEX idx_base_build_items_base ON base_build_items(base_id);
CREATE INDEX idx_base_build_items_base_status ON base_build_items(base_id, status);
CREATE INDEX idx_base_build_items_due ON base_build_items(((in_prod_data->>'completion_date')::bigint)) WHERE status = 'IN_PRODUCTION';

CREATE TABLE base_tech_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id             BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    prototype_id        BIGINT  NOT NULL REFERENCES tech_item_prototypes(id) ON DELETE RESTRICT,
    status              TEXT    NOT NULL,
    -- Tech in progress data: {"start_date": bigint, "completion_date": bigint, "crystals_skip_price": int}
    in_progress_data   JSONB,
    -- Tech done data: {"researched_at": bigint}
    done_data          JSONB,
    -- bookkeeping
    created_at          BIGINT  NOT NULL,
    -- constraints
    CONSTRAINT chk_tech_status CHECK (status IN ('IN_PROGRESS','DONE')),
    CONSTRAINT chk_tech_in_progress_json CHECK (
        status <> 'IN_PROGRESS' OR in_progress_data IS NOT NULL
    ),
    CONSTRAINT chk_tech_done_json CHECK (
        status <> 'DONE' OR done_data IS NOT NULL
    )
);
CREATE INDEX idx_base_tech_items_base ON base_tech_items(base_id);
CREATE INDEX idx_base_tech_items_base_status ON base_tech_items(base_id, status);
CREATE INDEX idx_base_tech_items_due ON base_tech_items(((in_progress_data->>'completion_date')::bigint)) WHERE status = 'IN_PROGRESS';
CREATE UNIQUE INDEX uq_base_tech_items_done ON base_tech_items(base_id, prototype_id) WHERE status = 'DONE';
CREATE UNIQUE INDEX uq_base_tech_items_in_progress ON base_tech_items(base_id, prototype_id) WHERE status = 'IN_PROGRESS';

CREATE TABLE base_storage_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id      BIGINT  NOT NULL REFERENCES user_bases(id) ON DELETE CASCADE,
    prototype_id BIGINT  NOT NULL REFERENCES storage_item_prototypes(id) ON DELETE RESTRICT,
    status       TEXT    NOT NULL,
    -- Storage present data: {"expires_at": bigint, "is_active": boolean}
    present_data    JSONB,
    -- Storage dynamic state: JSONB structure varies by prototype
    state        JSONB   NOT NULL DEFAULT '{}'::jsonb,
    created_at   BIGINT  NOT NULL,
    CONSTRAINT chk_storage_status CHECK (status IN ('PRESENT'))
);
CREATE INDEX idx_base_storage_items_base ON base_storage_items(base_id);
CREATE INDEX idx_base_storage_items_base_status ON base_storage_items(base_id, status);
CREATE INDEX idx_base_storage_items_expires ON base_storage_items(((present_data->>'expires_at')::bigint));
