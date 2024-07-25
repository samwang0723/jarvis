-- Drop the table
DROP TABLE IF EXISTS stocks;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_stocks_updated_at ON stocks CASCADE;
