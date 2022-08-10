-- +migrate Up notransaction
ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT UNIQUE NOT NULL;

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS email;
