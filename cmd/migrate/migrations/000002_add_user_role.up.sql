CREATE TYPE user_role AS ENUM ('teacher', 'student');

ALTER TABLE users
ADD COLUMN role user_role NOT NULL DEFAULT 'student';