CREATE TABLE daily_closes (
    id BIGINT NOT NULL PRIMARY KEY,
    stock_id VARCHAR(8) NOT NULL,
    exchange_date VARCHAR(32) NOT NULL,
    trade_shares BIGINT DEFAULT 0,
    transactions BIGINT DEFAULT 0,
    turnover BIGINT DEFAULT 0,
    open DECIMAL(8,2) NOT NULL,
    close DECIMAL(8,2) NOT NULL,
    high DECIMAL(8,2) NOT NULL,
    low DECIMAL(8,2) NOT NULL,
    price_diff DECIMAL(6,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_transactions ON daily_closes (trade_shares);
CREATE INDEX index_close ON daily_closes (close);
CREATE INDEX index_exchange_date ON daily_closes (exchange_date);
CREATE INDEX index_stock_id ON daily_closes (stock_id);

-- Create a trigger function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER update_updated_at BEFORE UPDATE ON daily_closes
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

