DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_unique_active_picked_stock_per_user;

DROP TRIGGER IF EXISTS update_users_updated_at ON users CASCADE;
