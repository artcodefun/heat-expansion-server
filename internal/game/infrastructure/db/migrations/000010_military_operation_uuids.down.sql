DROP INDEX IF EXISTS game.idx_military_operations_operation_uuid;

ALTER TABLE game.military_operations
    DROP COLUMN IF EXISTS operation_uuid;