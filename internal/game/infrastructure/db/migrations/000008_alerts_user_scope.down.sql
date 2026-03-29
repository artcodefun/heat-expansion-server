DROP INDEX IF EXISTS game.idx_alerts_user_id;

ALTER TABLE game.alerts
DROP CONSTRAINT IF EXISTS alerts_user_id_fkey;

ALTER TABLE game.alerts
DROP COLUMN IF EXISTS user_id;
