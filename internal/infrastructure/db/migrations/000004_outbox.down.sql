-- Down migration for domain events outbox table

DROP INDEX IF EXISTS idx_domain_events_published_id;
DROP TABLE IF EXISTS domain_events;
