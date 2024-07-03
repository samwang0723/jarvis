-- Drop the trigger
DROP TRIGGER IF EXISTS update_updated_at ON daily_closes;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the indexes
DROP INDEX IF EXISTS index_transactions;
DROP INDEX IF EXISTS index_close;
DROP INDEX IF EXISTS index_exchange_date;
DROP INDEX IF EXISTS index_stock_id;

-- Drop the table
DROP TABLE IF EXISTS daily_closes;

