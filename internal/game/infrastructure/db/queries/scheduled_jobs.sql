-- Scheduled jobs queries for durable scheduler

-- name: InsertScheduledJob :one
INSERT INTO game.scheduled_jobs (
    kind, payload, execute_at, created_at, dispatched
) VALUES (
    @kind, @payload, @execute_at, @created_at, FALSE
)
RETURNING id;

-- name: ClaimDueScheduledJobs :many
SELECT id, kind, payload, execute_at, created_at, dispatched, dispatched_at
FROM game.scheduled_jobs
WHERE dispatched = FALSE
  AND execute_at <= $1
ORDER BY execute_at ASC, id ASC
FOR UPDATE SKIP LOCKED
LIMIT $2;

-- name: MarkScheduledJobDispatched :exec
UPDATE game.scheduled_jobs
SET dispatched = TRUE,
    dispatched_at = @dispatched_at
WHERE id = @id;

-- name: GetNextScheduledJob :one
SELECT id, kind, payload, execute_at, created_at, dispatched, dispatched_at
FROM game.scheduled_jobs
WHERE dispatched = FALSE
ORDER BY execute_at ASC, id ASC
LIMIT 1;

