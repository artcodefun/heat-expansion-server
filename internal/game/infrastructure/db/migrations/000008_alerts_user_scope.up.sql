ALTER TABLE game.alerts
ADD COLUMN user_id UUID;

UPDATE game.alerts a
SET user_id = b.user_id
FROM game.user_bases b
WHERE b.id = a.base_id;

ALTER TABLE game.alerts
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE game.alerts
ADD CONSTRAINT alerts_user_id_fkey
FOREIGN KEY (user_id) REFERENCES game.users(id) ON DELETE CASCADE;

CREATE INDEX idx_alerts_user_id ON game.alerts(user_id);
