-- +goose Up
ALTER TABLE users
ADD COLUMN balance INT NOT NULL DEFAULT 1000;

-- +goose Down
ALTER TABLES users
DROP COLUMN balance;