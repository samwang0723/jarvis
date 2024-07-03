CREATE INDEX idx_stock_id_exchange_date_desc ON daily_closes (stock_id, exchange_date DESC);
CREATE INDEX idx_covering ON daily_closes (stock_id, exchange_date, close);
CREATE INDEX idx_covering_high ON daily_closes (stock_id, exchange_date, high);
