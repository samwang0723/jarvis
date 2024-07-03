CREATE TABLE stake_concentration (
    id BIGINT NOT NULL PRIMARY KEY,
    stock_id VARCHAR(8) NOT NULL,
    exchange_date VARCHAR(32) NOT NULL,
    sum_buy_shares BIGINT DEFAULT 0,
    sum_sell_shares BIGINT DEFAULT 0,
    avg_buy_price DECIMAL(8,2) NOT NULL,
    avg_sell_price DECIMAL(8,2) NOT NULL,
    concentration_1 DECIMAL(3,2) DEFAULT 0.0,
    concentration_5 DECIMAL(3,2) DEFAULT 0.0,
    concentration_10 DECIMAL(3,2) DEFAULT 0.0,
    concentration_20 DECIMAL(3,2) DEFAULT 0.0,
    concentration_60 DECIMAL(3,2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_exchange_date ON stake_concentration (exchange_date);
CREATE INDEX index_stock_id ON stake_concentration (stock_id);

-- Create a trigger function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER update_updated_at BEFORE UPDATE ON stake_concentration
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

