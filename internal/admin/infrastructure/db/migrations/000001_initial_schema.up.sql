CREATE SCHEMA IF NOT EXISTS admin;

CREATE TABLE admin.admins (
    id                UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    username          TEXT    NOT NULL UNIQUE,
    password_hash     TEXT,
    invite_token      TEXT,
    active            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        BIGINT  NOT NULL,
    updated_at        BIGINT  NOT NULL
);

CREATE TABLE admin.sessions (
    token      TEXT   PRIMARY KEY,
    admin_id   UUID   NOT NULL REFERENCES admin.admins (id),
    expires_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX admin_sessions_admin_id_idx   ON admin.sessions (admin_id);
CREATE INDEX admin_sessions_expires_at_idx ON admin.sessions (expires_at);
