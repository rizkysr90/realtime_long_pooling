-- migrate:up
CREATE TYPE order_status AS ENUM ('new', 'accepted', 'preparing', 'ready', 'picked_up', 'delivered');

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL,
    status order_status NOT NULL DEFAULT 'new',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

