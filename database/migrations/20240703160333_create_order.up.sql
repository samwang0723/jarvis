CREATE TABLE orders (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    stock_id VARCHAR(8) NOT NULL,
    buy_price DECIMAL(8, 2) NOT NULL,
    buy_quantity BIGINT NOT NULL DEFAULT 0,
    buy_exchange_date VARCHAR(32) NOT NULL,
    sell_price DECIMAL(8, 2) NOT NULL,
    sell_quantity BIGINT NOT NULL DEFAULT 0,
    sell_exchange_date VARCHAR(32) NOT NULL,
    profitable_price DECIMAL(8, 2) NOT NULL,
    status VARCHAR(32) NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX index_user_id ON orders(user_id);
CREATE INDEX index_user_stock_id ON orders(user_id, stock_id);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function before update
CREATE TRIGGER update_updated_at_trigger
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
