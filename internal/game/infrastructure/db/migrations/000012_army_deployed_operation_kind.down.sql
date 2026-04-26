-- Remove operation_kind from deployed army payloads and restore operation_id-only index.

DROP INDEX IF EXISTS idx_base_army_items_operation_kind_identity;
DROP INDEX IF EXISTS idx_base_army_items_operation_identity;

CREATE INDEX idx_base_army_items_operation
ON game.base_army_items (((deployed_data->>'operation_id')::bigint))
WHERE status = 'DEPLOYED';

UPDATE game.base_army_items
SET deployed_data = (deployed_data - 'operation_kind') - 'operation_type'
WHERE status = 'DEPLOYED'
  AND deployed_data IS NOT NULL
  AND ((deployed_data ? 'operation_kind') OR (deployed_data ? 'operation_type'));
