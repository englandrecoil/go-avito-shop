-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, username, hashed_password, balance)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    1000
)
RETURNING id, created_at, updated_at, username, balance;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;