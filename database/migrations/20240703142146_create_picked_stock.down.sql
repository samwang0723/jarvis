DROP TABLE IF EXISTS picked_stocks;
DROP INDEX IF EXISTS unique_active_picked_stock_per_user;
DROP TRIGGER IF EXISTS update_picked_stocks_updated_at ON picked_stocks CASCADE;
