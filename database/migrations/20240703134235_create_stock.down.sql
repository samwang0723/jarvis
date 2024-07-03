-- Drop the trigger
DROP TRIGGER IF EXISTS update_updated_at ON stocks;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the table
DROP TABLE IF EXISTS stocks;
