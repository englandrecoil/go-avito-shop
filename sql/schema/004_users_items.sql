-- +goose Up
CREATE TABLE users_items(
    id UUID PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK(quantity > 0),
    UNIQUE (user_id, item_id)
);

-- +goose Down
DROP TABLE users_items;