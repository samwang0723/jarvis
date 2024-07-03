-- Drop the trigger
DROP TRIGGER IF EXISTS update_updated_at ON three_primary;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the indexes
DROP INDEX IF EXISTS index_foreign;
DROP INDEX IF EXISTS index_trust;
DROP INDEX IF EXISTS index_dealer;
DROP INDEX IF EXISTS index_hedging;
DROP INDEX IF EXISTS index_exchange_date;
DROP INDEX IF EXISTS index_stock_id;

-- Drop the table
DROP TABLE IF EXISTS three_primary;
