-- Down migration: drops all tables in reverse dependency order (indexes are dropped automatically with their tables)

DROP TABLE game.activities;

DROP TABLE game.scan_reports;

-- Drop location-related child tables before parents
DROP TABLE game.dangerous_locations;
DROP TABLE game.resource_locations;

DROP TABLE game.military_operations;

DROP TABLE game.user_bases;

DROP TABLE game.sectors;

DROP TABLE game.users;

DROP SCHEMA game CASCADE;
