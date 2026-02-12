-- Scheduled jobs table for durable scheduler (PostgreSQL)

CREATE TABLE game.scheduled_jobs (
    id            BIGSERIAL PRIMARY KEY,
    kind          TEXT   NOT NULL,
    payload       JSONB  NOT NULL,
    execute_at    BIGINT NOT NULL,
    created_at    BIGINT NOT NULL,
    dispatched    BOOLEAN NOT NULL DEFAULT FALSE,
    dispatched_at BIGINT
);

-- Index to efficiently claim due, undispatched jobs by time and id
CREATE INDEX idx_scheduled_jobs_dispatched_execute_at_id
    ON game.scheduled_jobs (dispatched, execute_at, id);
