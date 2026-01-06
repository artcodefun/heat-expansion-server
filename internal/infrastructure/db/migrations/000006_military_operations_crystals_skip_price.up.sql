-- Add crystals_skip_price column to military_operations for operation speed-up pricing

ALTER TABLE military_operations
    ADD COLUMN crystals_skip_price INTEGER NOT NULL DEFAULT 0;
