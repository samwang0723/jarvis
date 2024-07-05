-- name: CreateDailyClose :exec
INSERT INTO daily_closes (
    stock_id, exchange_date, trade_shares, transactions, turnover, open, close, high, low, price_diff
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: BatchUpsertDailyClose :exec
INSERT INTO daily_closes (
    stock_id, exchange_date, trade_shares, transactions, turnover, open, close, high, low, price_diff
  ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
  )
ON CONFLICT (stock_id, exchange_date) DO UPDATE SET
    trade_shares = EXCLUDED.trade_shares,
    transactions = EXCLUDED.transactions,
    turnover = EXCLUDED.turnover,
    open = EXCLUDED.open,
    close = EXCLUDED.close,
    high = EXCLUDED.high,
    low = EXCLUDED.low,
    price_diff = EXCLUDED.price_diff;

-- name: HasDailyClose :one
SELECT stock_id FROM daily_closes
WHERE exchange_date = $1
LIMIT 1;

-- name: ListDailyClose :many
SELECT t.id, t.stock_id, t.exchange_date, t.transactions,
       FLOOR(t.trade_shares/1000) AS trade_shares, FLOOR(t.turnover/1000) AS turnover,
       t.open, t.high, t.close, t.low, t.price_diff, t.created_at, t.updated_at, t.deleted_at
FROM (
    SELECT daily_closes.id FROM daily_closes
    WHERE daily_closes.exchange_date >= @start_date AND daily_closes.stock_id = @stock_id
    AND (@end_date::text = '' OR daily_closes.exchange_date <= @end_date::tex)
    ORDER BY daily_closes.exchange_date DESC
    LIMIT $1 OFFSET $2
) q
JOIN daily_closes t ON t.id = q.id;

-- name: ListLatestPrice :many
WITH RankedCloses AS (
    SELECT 
        daily_closes.stock_id,
        daily_closes.close,
        daily_closes.exchange_date,
        ROW_NUMBER() OVER (PARTITION BY daily_closes.stock_id ORDER BY daily_closes.exchange_date DESC) AS rn
    FROM 
        daily_closes
    WHERE 
        daily_closes.stock_id = ANY(@stock_ids::text[])
)
SELECT 
    RankedCloses.stock_id,
    RankedCloses.close,
    RankedCloses.exchange_date
FROM 
    RankedCloses
WHERE 
    rn = 1;

-- name: CountDailyClose :one
SELECT COUNT(daily_closes.*)
FROM daily_closes
WHERE daily_closes.stock_id = @stock_id AND daily_closes.exchange_date >= @start_date
AND (@end_date::text = '' OR daily_closes.exchange_date <= @end_date::text);
