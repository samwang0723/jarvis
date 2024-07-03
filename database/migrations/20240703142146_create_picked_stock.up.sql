CREATE TABLE picked_stocks (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    stock_id VARCHAR(8) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (user_id, stock_id, deleted_at)
);

-- Create a trigger to update the updated_at column on row update
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_picked_stocks_updated_at
BEFORE UPDATE ON picked_stocks
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
