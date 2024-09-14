
ALTER TABLE files ADD COLUMN expiration_date TIMESTAMP;

-- Update existing files with a default expiration date (e.g., 30 days from now)
UPDATE files SET expiration_date = created_at + INTERVAL '30 days' WHERE expiration_date IS NULL;

-- Make expiration_date NOT NULL after updating existing records
ALTER TABLE files ALTER COLUMN expiration_date SET NOT NULL;

ALTER TABLE files
ALTER COLUMN expiration_date SET DEFAULT (CURRENT_TIMESTAMP + INTERVAL '30 days'); 