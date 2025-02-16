-- name: DeductBalance :exec
UPDATE users
SET balance = balance - $2  
WHERE id = $1;

-- name: AddBalance :exec
UPDATE users
SET balance = balance + $2
WHERE id = $1;

-- name: InsertTransaction :exec
INSERT INTO transactions (sender_id, receiver_id, amount)
VALUES ($1, $2, $3);
