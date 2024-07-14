-- Drop the table
DROP TABLE IF EXISTS three_primary;

-- Drop the indexes
DROP INDEX IF EXISTS idx_three_primary_id_exchange_date_desc;
DROP INDEX IF EXISTS idx_three_primary_stock_date_desc_full;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_three_primary_updated_at ON three_primary CASCADE;
