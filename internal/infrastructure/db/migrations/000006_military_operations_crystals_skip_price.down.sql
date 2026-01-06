-- Drop crystals_skip_price column from military_operations

ALTER TABLE military_operations
    DROP COLUMN crystals_skip_price;
