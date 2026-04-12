DROP INDEX IF EXISTS game.idx_diplomatic_messages_receiver_sender_unread;
DROP INDEX IF EXISTS game.idx_diplomatic_messages_request_created;
DROP INDEX IF EXISTS game.idx_diplomatic_messages_sender_created;
DROP INDEX IF EXISTS game.idx_diplomatic_messages_receiver_created;
DROP TABLE IF EXISTS game.diplomatic_messages;

DROP INDEX IF EXISTS game.idx_diplomatic_requests_pair_kind_pending;
DROP INDEX IF EXISTS game.idx_diplomatic_requests_sender_created;
DROP INDEX IF EXISTS game.idx_diplomatic_requests_receiver_status_created;
DROP TABLE IF EXISTS game.diplomatic_requests;

DROP INDEX IF EXISTS game.idx_diplomatic_relationships_user_b;
DROP INDEX IF EXISTS game.idx_diplomatic_relationships_user_a;
DROP TABLE IF EXISTS game.diplomatic_relationships;

ALTER TABLE game.alerts
ALTER COLUMN base_id SET NOT NULL;
