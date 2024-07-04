BEGIN;

CREATE TABLE stake_concentration (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    sum_buy_shares bigint DEFAULT 0,
    sum_sell_shares bigint DEFAULT 0,
    avg_buy_price numeric(8,2) NOT NULL,
    avg_sell_price numeric(8,2) NOT NULL,
    concentration_1 numeric(5,2) DEFAULT 0.0,
    concentration_5 numeric(5,2) DEFAULT 0.0,
    concentration_10 numeric(5,2) DEFAULT 0.0,
    concentration_20 numeric(5,2) DEFAULT 0.0,
    concentration_60 numeric(5,2) DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_stake_concentration_exchange_date ON stake_concentration (exchange_date);
CREATE INDEX index_stake_concentration_stock_id ON stake_concentration (stock_id);

-- Create the trigger
CREATE TRIGGER update_stake_concentration_updated_at BEFORE UPDATE ON stake_concentration
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
