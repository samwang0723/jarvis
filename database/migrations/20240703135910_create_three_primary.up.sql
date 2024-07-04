BEGIN;

CREATE TABLE three_primary (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    foreign_trade_shares bigint DEFAULT 0,
    trust_trade_shares bigint DEFAULT 0,
    dealer_trade_shares bigint DEFAULT 0,
    hedging_trade_shares bigint DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_three_primary_foreign ON three_primary (foreign_trade_shares);
CREATE INDEX index_three_primary_trust ON three_primary (trust_trade_shares);
CREATE INDEX index_three_primary_dealer ON three_primary (dealer_trade_shares);
CREATE INDEX index_three_primary_hedging ON three_primary (hedging_trade_shares);
CREATE INDEX index_three_primary_exchange_date ON three_primary (exchange_date);
CREATE INDEX index_three_primary_stock_id ON three_primary (stock_id);

-- Create the trigger
CREATE TRIGGER update_three_primary_updated_at BEFORE UPDATE ON three_primary
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
