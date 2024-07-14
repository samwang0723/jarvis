BEGIN;

CREATE TABLE orders (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    stock_id varchar(8) NOT NULL,
    buy_price numeric(8, 2) NOT NULL,
    buy_quantity bigint NOT NULL DEFAULT 0,
    buy_exchange_date varchar(32) NOT NULL,
    sell_price numeric(8, 2) NOT NULL,
    sell_quantity bigint NOT NULL DEFAULT 0,
    sell_exchange_date varchar(32) NOT NULL,
    profitable_price numeric(8, 2) NOT NULL,
    status varchar(32) NOT NULL,
    version integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_stock_id FOREIGN KEY (stock_id) REFERENCES stocks (id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create indexes
CREATE INDEX idx_order_user_id ON orders(user_id);
CREATE INDEX idx_order_user_stock_id ON orders(user_id, stock_id);

-- Create a trigger to call the function before update
CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
