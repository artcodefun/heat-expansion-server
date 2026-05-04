CREATE TABLE auth.password_reset_tokens (
    id          UUID    PRIMARY KEY,
    account_id  UUID    NOT NULL REFERENCES auth.users(id),
    token_hash  TEXT    NOT NULL,
    expires_at  BIGINT  NOT NULL,
    used_at     BIGINT
);

CREATE INDEX idx_auth_prt_token_hash ON auth.password_reset_tokens (token_hash);
