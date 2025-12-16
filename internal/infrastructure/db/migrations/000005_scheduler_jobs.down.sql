-- Down migration for scheduled jobs table

DROP INDEX IF EXISTS idx_scheduled_jobs_dispatched_execute_at_id;
DROP TABLE IF EXISTS scheduled_jobs;
