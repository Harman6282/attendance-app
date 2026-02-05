ALTER TABLE users
DROP COLUMN IF EXISTS role;

-- remove enum type
DROP TYPE IF EXISTS user_role;