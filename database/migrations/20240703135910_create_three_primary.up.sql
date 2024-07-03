CREATE TABLE three_primary (
    id BIGINT NOT NULL PRIMARY KEY,
    stock_id VARCHAR(8) NOT NULL,
    exchange_date VARCHAR(32) NOT NULL,
    foreign_trade_shares BIGINT DEFAULT 0,
    trust_trade_shares BIGINT DEFAULT 0,
    dealer_trade_shares BIGINT DEFAULT 0,
    hedging_trade_shares BIGINT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_foreign ON three_primary (foreign_trade_shares);
CREATE INDEX index_trust ON three_primary (trust_trade_shares);
CREATE INDEX index_dealer ON three_primary (dealer_trade_shares);
CREATE INDEX index_hedging ON three_primary (hedging_trade_shares);
CREATE INDEX index_exchange_date ON three_primary (exchange_date);
CREATE INDEX index_stock_id ON three_primary (stock_id);

-- Create a trigger function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER update_updated_at BEFORE UPDATE ON three_primary
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

