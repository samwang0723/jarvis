DROP TABLE IF EXISTS transactions;

DROP INDEX IF EXISTS idx_transaction_user_id;
DROP INDEX IF EXISTS idx_transaction_order_id;

DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions CASCADE;

