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
CREATE INDEX idx_three_primary_id_exchange_date_desc ON three_primary (id, exchange_date DESC);
CREATE INDEX idx_three_primary_stock_date_desc_full ON three_primary
(stock_id, exchange_date DESC, foreign_trade_shares, trust_trade_shares, dealer_trade_shares, hedging_trade_shares);

-- Create the trigger
CREATE TRIGGER update_three_primary_updated_at BEFORE UPDATE ON three_primary
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
