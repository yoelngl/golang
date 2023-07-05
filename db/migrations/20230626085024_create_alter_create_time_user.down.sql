ALTER TABLE users RENAME COLUMN created_at TO created_time;
ALTER TABLE users RENAME COLUMN updated_at TO  updated_time;
ALTER TABLE users MODIFY COLUMN created_at TIMESTAMP NOT NULL;
