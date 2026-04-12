ALTER TABLE game.alerts
ALTER COLUMN base_id DROP NOT NULL;

CREATE TABLE game.diplomatic_relationships (
    id UUID PRIMARY KEY,
    user_a_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    user_b_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    changed_by_user_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    changed_at BIGINT NOT NULL,
    war_declared_at BIGINT,
    war_attacks_allowed_at BIGINT,
    neutrality_protected_until BIGINT,
    CONSTRAINT diplomatic_relationships_pair_unique UNIQUE (user_a_id, user_b_id),
    CONSTRAINT diplomatic_relationships_distinct_users CHECK (user_a_id <> user_b_id)
);

CREATE INDEX idx_diplomatic_relationships_user_a ON game.diplomatic_relationships(user_a_id);
CREATE INDEX idx_diplomatic_relationships_user_b ON game.diplomatic_relationships(user_b_id);

CREATE TABLE game.diplomatic_requests (
    id UUID PRIMARY KEY,
    sender_user_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    receiver_user_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    sender_base_id BIGINT REFERENCES game.user_bases(id) ON DELETE SET NULL,
    receiver_base_id BIGINT REFERENCES game.user_bases(id) ON DELETE SET NULL,
    kind TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    resolved_at BIGINT,
    expires_at BIGINT NOT NULL
);

CREATE INDEX idx_diplomatic_requests_receiver_status_created ON game.diplomatic_requests(receiver_user_id, status, created_at DESC);
CREATE INDEX idx_diplomatic_requests_sender_created ON game.diplomatic_requests(sender_user_id, created_at DESC);
CREATE INDEX idx_diplomatic_requests_pair_kind_pending ON game.diplomatic_requests(sender_user_id, receiver_user_id, kind) WHERE status = 'PENDING';

CREATE TABLE game.diplomatic_messages (
    id UUID PRIMARY KEY,
    sender_user_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    receiver_user_id UUID NOT NULL REFERENCES game.users(id) ON DELETE CASCADE,
    sender_base_id BIGINT REFERENCES game.user_bases(id) ON DELETE SET NULL,
    receiver_base_id BIGINT REFERENCES game.user_bases(id) ON DELETE SET NULL,
    request_id UUID REFERENCES game.diplomatic_requests(id) ON DELETE SET NULL,
    reply_to_message_id UUID REFERENCES game.diplomatic_messages(id) ON DELETE SET NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX idx_diplomatic_messages_receiver_created ON game.diplomatic_messages(receiver_user_id, created_at DESC);
CREATE INDEX idx_diplomatic_messages_sender_created ON game.diplomatic_messages(sender_user_id, created_at DESC);
CREATE INDEX idx_diplomatic_messages_request_created ON game.diplomatic_messages(request_id, created_at DESC) WHERE request_id IS NOT NULL;
CREATE INDEX idx_diplomatic_messages_receiver_sender_unread ON game.diplomatic_messages(receiver_user_id, sender_user_id, created_at DESC) WHERE is_read = FALSE;
