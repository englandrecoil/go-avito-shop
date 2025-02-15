-- name: GetItemByName :one
SELECT * FROM items
WHERE name = $1;

-- name: PurchaseItemByID :exec
WITH updated_balance AS (
    UPDATE users
    SET balance = balance - (SELECT price FROM items WHERE id = $2)
    WHERE id = $1
    AND balance >= (SELECT price FROM items WHERE id = $2)
    RETURNING id
)
INSERT INTO users_items (id, user_id, item_id, quantity)
VALUES (gen_random_uuid(), $1, $2, 1)
ON CONFLICT (user_id, item_id) 
DO UPDATE SET quantity = users_items.quantity + EXCLUDED.quantity
WHERE EXISTS (SELECT 1 FROM updated_balance);
