-- Add stable UUIDs to military operations so external provenance can stop depending on the numeric PK.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE game.military_operations
    ADD COLUMN operation_uuid UUID;

UPDATE game.military_operations
SET operation_uuid = gen_random_uuid()
WHERE operation_uuid IS NULL;

ALTER TABLE game.military_operations
    ALTER COLUMN operation_uuid SET NOT NULL;

CREATE UNIQUE INDEX idx_military_operations_operation_uuid ON game.military_operations(operation_uuid);