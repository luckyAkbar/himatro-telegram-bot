-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "users" (
    id BIGINT PRIMARY KEY,
    user_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "sessions" (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    access_token TEXT UNIQUE NOT NULL,
    refresh_token TEXT UNIQUE NOT NULL,
    expired_at TIMESTAMP NOT NULL
);

ALTER TABLE "sessions" ADD CONSTRAINT "sessions_user_id" FOREIGN KEY (user_id) REFERENCES "users" (id);

-- +migrate Down
DROP TABLE IF EXISTS "sessions";