-- Add DEPLOYED lifecycle support for base storage items used by trade operations.

ALTER TABLE game.base_storage_items
ADD COLUMN IF NOT EXISTS deployed_data JSONB;

ALTER TABLE game.base_storage_items
DROP CONSTRAINT IF EXISTS chk_storage_status;

ALTER TABLE game.base_storage_items
ADD CONSTRAINT chk_storage_status CHECK (status IN ('PRESENT', 'DEPLOYED'));

ALTER TABLE game.base_storage_items
ADD CONSTRAINT chk_storage_present_json CHECK (
    status <> 'PRESENT' OR present_data IS NOT NULL
);

ALTER TABLE game.base_storage_items
ADD CONSTRAINT chk_storage_deployed_json CHECK (
    status <> 'DEPLOYED' OR deployed_data IS NOT NULL
);

DROP INDEX IF EXISTS idx_base_storage_items_operation_identity;
CREATE INDEX idx_base_storage_items_operation_identity
ON game.base_storage_items (
    (deployed_data->>'operation_kind'),
    ((deployed_data->>'operation_id')::bigint)
)
WHERE status = 'DEPLOYED';
