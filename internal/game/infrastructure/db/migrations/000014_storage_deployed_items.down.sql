-- Remove DEPLOYED lifecycle support for base storage items.

DROP INDEX IF EXISTS idx_base_storage_items_operation_identity;

ALTER TABLE game.base_storage_items
DROP CONSTRAINT IF EXISTS chk_storage_deployed_json;

ALTER TABLE game.base_storage_items
DROP CONSTRAINT IF EXISTS chk_storage_present_json;

ALTER TABLE game.base_storage_items
DROP CONSTRAINT IF EXISTS chk_storage_status;

ALTER TABLE game.base_storage_items
ADD CONSTRAINT chk_storage_status CHECK (status IN ('PRESENT'));

ALTER TABLE game.base_storage_items
DROP COLUMN IF EXISTS deployed_data;
