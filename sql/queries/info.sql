-- name: GetReceivedHistory :many
SELECT users.username AS sender_name, amount AS received
FROM transactions 
INNER JOIN users on users.id = transactions.sender_id
WHERE receiver_id = $1;


-- name: GetSentHistory :many
SELECT users.username AS receiver_name, amount AS sent
FROM transactions
INNER JOIN users on users.id = transactions.receiver_id
WHERE sender_id = $1;

-- name: GetInventory :many
SELECT items.name AS item_name, users_items.quantity
FROM users_items
INNER JOIN items on items.id = users_items.item_id
INNER JOIN users on users.id = users_items.user_id
WHERE user_id = $1
ORDER BY quantity DESC;
