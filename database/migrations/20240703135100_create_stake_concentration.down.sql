-- Drop the table
DROP TABLE IF EXISTS stake_concentration;

-- Drop the indexes
DROP INDEX IF EXISTS index_stake_concentration_exchange_date;
DROP INDEX IF EXISTS index_stake_concentration_stock_id;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_stake_concentration_updated_at ON stake_concentration CASCADE;
