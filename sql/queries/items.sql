-- name: GetItemByName :one
SELECT * FROM items
WHERE name = $1;

-- name: PurchaseItemByID :one
WITH item_price AS (
    SELECT price FROM items WHERE id = $2
),
updated_balance AS (
    UPDATE users
    SET balance = balance - (SELECT price FROM item_price)
    WHERE id = $1 AND balance >= (SELECT price FROM item_price)
    RETURNING id
)
INSERT INTO users_items (id, user_id, item_id, quantity)
SELECT gen_random_uuid(), $1, $2, 1
FROM updated_balance
ON CONFLICT (user_id, item_id) 
DO UPDATE SET quantity = users_items.quantity + EXCLUDED.quantity
RETURNING *;
