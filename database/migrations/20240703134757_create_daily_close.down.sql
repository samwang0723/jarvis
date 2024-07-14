-- Drop the table
DROP TABLE IF EXISTS daily_closes;

-- Drop the indexes
drop index if exists idx_daily_closes_exchange_date;
drop index if exists idx_daily_closes_stock_id;
drop index if exists idx_daily_closes_covering_high;
drop index if exists idx_daily_closes_stock_date_desc_full;

-- Drop the trigger
DROP TRIGGER IF EXISTS update_daily_closes_updated_at ON daily_closes CASCADE;
