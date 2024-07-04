BEGIN;

CREATE TABLE daily_closes (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    trade_shares bigint DEFAULT 0,
    transactions bigint DEFAULT 0,
    turnover bigint DEFAULT 0,
    open numeric(8,2) NOT NULL,
    close numeric(8,2) NOT NULL,
    high numeric(8,2) NOT NULL,
    low numeric(8,2) NOT NULL,
    price_diff numeric(6,2) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_timestamp,
    updated_at timestamp NOT NULL DEFAULT CURRENT_timestamp,
    deleted_at timestamp NULL,
    UNIQUE (stock_id, exchange_date)
);

-- Create indexes
CREATE INDEX index_daily_closes_transactions ON daily_closes (trade_shares);
CREATE INDEX index_daily_closes_close ON daily_closes (close);
CREATE INDEX index_daily_closes_exchange_date ON daily_closes (exchange_date);
CREATE INDEX index_daily_closes_stock_id ON daily_closes (stock_id);
CREATE INDEX index_daily_closes_stock_id_exchange_date_desc ON daily_closes (stock_id, exchange_date DESC);
CREATE INDEX index_daily_closes_covering ON daily_closes (stock_id, exchange_date, close);
CREATE INDEX index_daily_closes_covering_high ON daily_closes (stock_id, exchange_date, high);

-- Create the trigger
CREATE TRIGGER update_daily_closes_updated_at BEFORE UPDATE ON daily_closes
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

