DROP TABLE IF EXISTS orders;

DROP INDEX IF EXISTS idx_order_user_id;
DROP INDEX IF EXISTS idx_order_user_stock_id;

DROP TRIGGER IF EXISTS update_orders_updated_at ON orders CASCADE;

