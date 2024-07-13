-- name: GetStakeConcentrationByStockID :one
SELECT 
  id,
  stock_id, 
  exchange_date, 
  sum_buy_shares, 
  sum_sell_shares, 
  avg_buy_price, 
  avg_sell_price, 
  concentration_1::numeric, 
  concentration_5::numeric, 
  concentration_10::numeric, 
  concentration_20::numeric, 
  concentration_60::numeric,
  created_at,
  updated_at,
  deleted_at
FROM stake_concentration
WHERE stock_id = $1 AND exchange_date = $2
LIMIT 1;

-- name: BatchUpsertStakeConcentration :exec
INSERT INTO stake_concentration (
  stock_id, 
  exchange_date, 
  sum_buy_shares, 
  sum_sell_shares, 
  avg_buy_price, 
  avg_sell_price, 
  concentration_1, 
  concentration_5, 
  concentration_10, 
  concentration_20, 
  concentration_60)
VALUES (
  unnest(@stock_id::varchar[]), 
  unnest(@exchange_date::varchar[]), 
  unnest(@sum_buy_shares::bigint[]), 
  unnest(@sum_sell_shares::bigint[]), 
  unnest(@avg_buy_price::numeric[]),
  unnest(@avg_sell_price::numeric[]),
  unnest(@concentration_1::numeric[]),
  unnest(@concentration_5::numeric[]),
  unnest(@concentration_10::numeric[]),
  unnest(@concentration_20::numeric[]),
  unnest(@concentration_60::numeric[])
)
ON CONFLICT (stock_id, exchange_date) DO UPDATE
SET sum_buy_shares = EXCLUDED.sum_buy_shares,
    sum_sell_shares = EXCLUDED.sum_sell_shares,
    avg_buy_price = EXCLUDED.avg_buy_price,
    avg_sell_price = EXCLUDED.avg_sell_price,
    concentration_1 = EXCLUDED.concentration_1,
    concentration_5 = EXCLUDED.concentration_5,
    concentration_10 = EXCLUDED.concentration_10,
    concentration_20 = EXCLUDED.concentration_20,
    concentration_60 = EXCLUDED.concentration_60;

-- name: GetStakeConcentrationsWithVolumes :many
SELECT a.trade_shares,
  COALESCE(b.sum_buy_shares, 0)::bigint - COALESCE(b.sum_sell_shares, 0)::bigint AS diff,
  a.exchange_date
FROM daily_closes a
LEFT JOIN stake_concentration b ON (a.stock_id, a.exchange_date) = (b.stock_id, b.exchange_date)
WHERE a.stock_id = $1 
AND a.exchange_date <= $2
ORDER BY a.exchange_date DESC
LIMIT 60;

-- name: HasStakeConcentration :one
SELECT EXISTS (
  SELECT 1 FROM stake_concentration
  WHERE exchange_date = $1
);

-- name: GetStakeConcentrationLatestDataPoint :one
SELECT exchange_date FROM stake_concentration
ORDER BY exchange_date DESC LIMIT 1;
