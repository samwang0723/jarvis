DROP TRIGGER IF EXISTS update_updated_at_trigger ON orders;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP INDEX IF EXISTS index_user_id;
DROP INDEX IF EXISTS index_user_stock_id;
DROP TABLE IF EXISTS orders;
