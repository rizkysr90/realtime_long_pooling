-- migrate:up
CREATE TYPE order_status AS ENUM ('new', 'accepted', 'preparing', 'ready', 'picked_up', 'delivered');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL,
    status order_status NOT NULL DEFAULT 'new',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE notification_type AS ENUM ('new_order', 'order_cancelled');

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    type notification_type NOT NULL,
    read_status BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- migrate:down

