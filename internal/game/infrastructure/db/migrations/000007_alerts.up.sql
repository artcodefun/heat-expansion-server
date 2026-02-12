CREATE TABLE game.alerts (
    id UUID PRIMARY KEY,
    base_id BIGINT NOT NULL REFERENCES game.user_bases(id) ON DELETE CASCADE,
    activity_id UUID REFERENCES game.activities(id) ON DELETE SET NULL,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL,
    expires_at BIGINT NOT NULL
);

CREATE INDEX idx_alerts_base_id ON game.alerts(base_id);
CREATE INDEX idx_alerts_expires_at ON game.alerts(expires_at);
