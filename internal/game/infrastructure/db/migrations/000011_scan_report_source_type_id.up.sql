ALTER TABLE game.scan_reports
    ADD COLUMN source_type TEXT NOT NULL DEFAULT 'UNKNOWN',
    ADD COLUMN source_id UUID;

UPDATE game.scan_reports sr
SET source_type = CASE
        WHEN sr.source_intel_item_id IS NOT NULL THEN 'INTEL'
        WHEN sr.source_scanner_id IS NOT NULL THEN 'SCANNER'
        WHEN sr.source_operation_id IS NOT NULL THEN 'OPERATION'
        ELSE 'UNKNOWN'
    END,
    source_id = CASE
        WHEN sr.source_intel_item_id IS NOT NULL THEN sr.source_intel_item_id
        WHEN sr.source_scanner_id IS NOT NULL THEN sr.source_scanner_id
        WHEN sr.source_operation_id IS NOT NULL THEN op.operation_uuid
        ELSE NULL
    END
FROM game.military_operations op
WHERE sr.source_operation_id = op.id;

UPDATE game.scan_reports
SET source_type = 'UNKNOWN',
    source_id = NULL
WHERE source_id IS NULL
  AND source_type = 'UNKNOWN';

ALTER TABLE game.scan_reports
        DROP COLUMN IF EXISTS source_operation_id,
        DROP COLUMN IF EXISTS source_scanner_id,
        DROP COLUMN IF EXISTS source_intel_item_id;

CREATE INDEX IF NOT EXISTS idx_scan_reports_source_type_id ON game.scan_reports(source_type, source_id);