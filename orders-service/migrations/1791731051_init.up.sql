CREATE TABLE orders (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    status      VARCHAR(20) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE order_items (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name     VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price    NUMERIC(12, 2) NOT NULL,
    PRIMARY KEY (order_id, name)
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);