-- Drop the table
DROP TABLE IF EXISTS daily_closes;

-- Drop the indexes
DROP INDEX IF EXISTS index_daily_closes_transactions;
DROP INDEX IF EXISTS index_daily_closes_close;
DROP INDEX IF EXISTS index_daily_closes_exchange_date;
DROP INDEX IF EXISTS index_daily_closes_stock_id;
DROP INDEX IF EXISTS index_daily_closes_stock_id_exchange_date_desc;
DROP INDEX IF EXISTS index_daily_closes_covering;
DROP INDEX IF EXISTS index_daily_closes_covering_high;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_daily_closes_updated_at ON daily_closes CASCADE;
