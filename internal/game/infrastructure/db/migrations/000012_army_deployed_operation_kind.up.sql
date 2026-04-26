-- Add operation_kind to deployed army payloads and switch deployed operation index to composite identity.

UPDATE game.base_army_items
SET deployed_data = CASE
  WHEN deployed_data ? 'operation_kind' THEN deployed_data
  WHEN deployed_data ? 'operation_type' THEN
    jsonb_set(
      deployed_data - 'operation_type',
      '{operation_kind}',
      deployed_data->'operation_type',
      true
    )
  ELSE
    jsonb_set(COALESCE(deployed_data, '{}'::jsonb), '{operation_kind}', to_jsonb('MILITARY'::text), true)
END
WHERE status = 'DEPLOYED'
  AND deployed_data IS NOT NULL;

DROP INDEX IF EXISTS idx_base_army_items_operation;
DROP INDEX IF EXISTS idx_base_army_items_operation_identity;
DROP INDEX IF EXISTS idx_base_army_items_operation_kind_identity;

CREATE INDEX idx_base_army_items_operation_kind_identity
ON game.base_army_items (
  (deployed_data->>'operation_kind'),
    ((deployed_data->>'operation_id')::bigint)
)
WHERE status = 'DEPLOYED';
