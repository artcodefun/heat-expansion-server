-- Down migration: drops all tables in reverse dependency order (indexes are dropped automatically with their tables)

DROP TABLE activities;

DROP TABLE scan_reports;

-- Drop location-related child tables before parents
DROP TABLE dangerous_locations;
DROP TABLE resource_locations;

DROP TABLE military_operations;

DROP TABLE user_bases;

DROP TABLE sectors;

DROP TABLE users;
