CREATE TABLE transactions (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    order_type VARCHAR(32) NOT NULL,
    credit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    debit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    status VARCHAR(32) NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT index_order_id FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- Create an index on order_id
CREATE INDEX index_order_id ON transactions(order_id);

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
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
