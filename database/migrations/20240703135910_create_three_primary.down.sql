-- Drop the table
DROP TABLE IF EXISTS three_primary;

-- Drop the indexes
DROP INDEX IF EXISTS index_three_primary_foreign;
DROP INDEX IF EXISTS index_three_primary_trust;
DROP INDEX IF EXISTS index_three_primary_dealer;
DROP INDEX IF EXISTS index_three_primary_hedging;
DROP INDEX IF EXISTS index_three_primary_exchange_date;
DROP INDEX IF EXISTS index_three_primary_stock_id;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_three_primary_updated_at ON three_primary CASCADE;
