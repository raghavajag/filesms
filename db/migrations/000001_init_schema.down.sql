-- Drop foreign key constraint before dropping tables
ALTER TABLE shared_file_urls DROP CONSTRAINT shared_file_urls_file_id_fkey;

-- Drop tables in the correct order
DROP TABLE IF EXISTS shared_file_urls;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS users;