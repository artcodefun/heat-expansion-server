DROP INDEX IF EXISTS game.idx_scan_reports_source_type_id;

ALTER TABLE game.scan_reports
        ADD COLUMN IF NOT EXISTS source_operation_id BIGINT REFERENCES game.military_operations(id) ON DELETE SET NULL,
        ADD COLUMN IF NOT EXISTS source_scanner_id UUID,
        ADD COLUMN IF NOT EXISTS source_intel_item_id UUID;

UPDATE game.scan_reports sr
SET source_operation_id = op.id
FROM game.military_operations op
WHERE sr.source_type = 'OPERATION'
    AND sr.source_id = op.operation_uuid;

UPDATE game.scan_reports
SET source_scanner_id = source_id
WHERE source_type = 'SCANNER'
    AND source_id IS NOT NULL;

UPDATE game.scan_reports
SET source_intel_item_id = source_id
WHERE source_type = 'INTEL'
    AND source_id IS NOT NULL;

ALTER TABLE game.scan_reports
    DROP COLUMN IF EXISTS source_id,
    DROP COLUMN IF EXISTS source_type;