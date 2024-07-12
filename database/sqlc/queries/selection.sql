-- name: LatestStockStatSnapshot :many
WITH latest AS (
  SELECT exchange_date 
  FROM stake_concentration
  ORDER BY exchange_date DESC 
  LIMIT 1
)
SELECT 
  s.stock_id, 
  c.name, 
  (c.category || '.' || c.market)::text AS category, 
  s.exchange_date, 
  d.open, 
  d.close, 
  d.high, 
  d.low, 
  d.price_diff,
  s.concentration_1, 
  s.concentration_5, 
  s.concentration_10, 
  s.concentration_20, 
  s.concentration_60, 
  floor(d.trade_shares/1000) AS volume, 
  floor(COALESCE(t.foreign_trade_shares,0)/1000) AS foreignc,
  floor(COALESCE(t.trust_trade_shares,0)/1000) AS trust, 
  floor(COALESCE(t.hedging_trade_shares,0)/1000) AS hedging,
  floor(COALESCE(t.dealer_trade_shares,0)/1000) AS dealer
FROM 
  stake_concentration s
LEFT JOIN latest ON 1=1
LEFT JOIN 
  stocks c ON c.id = s.stock_id
LEFT JOIN 
  daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = latest.exchange_date)
LEFT JOIN 
  three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = latest.exchange_date)
WHERE 
  (
    CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
    CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
    CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
    CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
    CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
  ) >= 4
  AND c.name IS NOT NULL
  AND s.exchange_date = latest.exchange_date
  AND d.trade_shares >= 1000000
ORDER BY 
  s.stock_id;

-- name: GetEligibleStocksFromDate :many
SELECT s.stock_id, c.market
FROM stake_concentration s
LEFT JOIN stocks c ON c.id = s.stock_id
LEFT JOIN daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = $1)
WHERE (
   CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
) >= 4
AND c.name IS NOT NULL
AND s.exchange_date = $1
AND d.trade_shares >= 1000000;

-- name: GetEligibleStocksFromPicked :many
SELECT p.stock_id, c.market
FROM picked_stocks p
LEFT JOIN stocks c ON c.id = p.stock_id 
WHERE p.deleted_at is null;

-- name: GetEligibleStocksFromOrder :many
SELECT o.stock_id, c.market
FROM orders o
LEFT JOIN stocks c ON c.id = o.stock_id 
WHERE o.status != 'closed';

-- name: ListSelectionsFromPicked :many
WITH latest AS (
  SELECT exchange_date 
  FROM stake_concentration
  ORDER BY exchange_date DESC 
  LIMIT 1
)
SELECT
  s.stock_id, 
  c.name, 
  (c.category || '.' || c.market)::text AS category, 
  s.exchange_date, 
  d.open, 
  d.close, 
  d.high, 
  d.low, 
  d.price_diff,
  s.concentration_1, 
  s.concentration_5, 
  s.concentration_10, 
  s.concentration_20, 
  s.concentration_60, 
  floor(d.trade_shares/1000) AS volume, 
  floor(COALESCE(t.foreign_trade_shares,0)/1000) AS foreignc,
  floor(COALESCE(t.trust_trade_shares,0)/1000) AS trust, 
  floor(COALESCE(t.hedging_trade_shares,0)/1000) AS hedging,
  floor(COALESCE(t.dealer_trade_shares,0)/1000) AS dealer
FROM stake_concentration s
LEFT JOIN latest ON 1=1
LEFT JOIN stocks c ON c.id = s.stock_id
LEFT JOIN daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = latest.exchange_date)
LEFT JOIN three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = latest.exchange_date)
WHERE s.stock_id = ANY(@stock_ids::text[])
AND c.name IS NOT NULL
AND s.exchange_date = latest.exchange_date
ORDER BY s.stock_id;

-- name: ListSelections :many
WITH average AS (
  SELECT 
    stock_id, 
    AVG(trade_shares) AS avg_volume
  FROM 
    daily_closes
  WHERE 
    exchange_date BETWEEN TO_CHAR(TO_DATE($1, 'YYYYMMDD') - INTERVAL '5' day, 'YYYYMMDD') 
      AND TO_CHAR(TO_DATE($1, 'YYYYMMDD') - INTERVAL '1' day, 'YYYYMMDD')
  GROUP BY 
    stock_id
)
SELECT 
  s.stock_id, 
  c.name, 
  (c.category || '.' || c.market)::text AS category, 
  s.exchange_date, 
  d.open, 
  d.close, 
  d.high, 
  d.low, 
  d.price_diff,
  s.concentration_1, 
  s.concentration_5, 
  s.concentration_10, 
  s.concentration_20, 
  s.concentration_60, 
  FLOOR(d.trade_shares/1000) AS volume, 
  FLOOR(COALESCE(t.foreign_trade_shares,0)/1000) AS foreignc,
  FLOOR(COALESCE(t.trust_trade_shares,0)/1000) AS trust, 
  FLOOR(COALESCE(t.hedging_trade_shares,0)/1000) AS hedging,
  FLOOR(COALESCE(t.dealer_trade_shares,0)/1000) AS dealer,
  a.avg_volume
FROM 
  stake_concentration s
LEFT JOIN stocks c ON c.id = s.stock_id
LEFT JOIN daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = $1)
LEFT JOIN three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = $1)
LEFT JOIN average a ON a.stock_id = s.stock_id
WHERE (
   CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
) >= 4
AND c.name IS NOT NULL
AND s.exchange_date = $1
AND d.trade_shares >= 3000000
AND a.avg_volume >= 1000000
ORDER BY s.stock_id;

-- name: GetStartDate :one
SELECT MIN(a.exchange_date)::text 
FROM (
  SELECT exchange_date 
  FROM stake_concentration
  WHERE exchange_date <= $1 
  GROUP BY exchange_date 
  ORDER BY exchange_date DESC 
  LIMIT 120
) AS a;

-- name: GetHighestPrice :many
SELECT stock_id, MAX(high)::numeric AS high 
FROM daily_closes 
WHERE exchange_date >= @start_date 
AND exchange_date < @end_date 
AND stock_id = ANY(@stock_ids::text[]) 
GROUP BY stock_id;

-- name: RetrieveDailyCloseHistoryWithDate :many
SELECT stock_id, exchange_date, close, trade_shares 
FROM daily_closes 
WHERE exchange_date >= @start_date 
AND exchange_date <= @end_date 
AND stock_id = ANY(@stock_ids::text[]) 
ORDER BY stock_id, exchange_date DESC;

-- name: RetrieveDailyCloseHistory :many
SELECT stock_id, exchange_date, close, trade_shares 
FROM daily_closes 
WHERE exchange_date >= @start_date 
AND exchange_date < @end_date
AND stock_id = ANY(@stock_ids::text[]) 
ORDER BY stock_id, exchange_date DESC;

-- name: RetrieveThreePrimaryHistoryWithDate :many
SELECT stock_id, exchange_date, 
floor(foreign_trade_shares/1000) AS foreign_trade_shares, 
floor(trust_trade_shares/1000) AS trust_trade_shares, 
floor(dealer_trade_shares/1000) AS dealer_trade_shares, 
floor(hedging_trade_shares/1000) AS hedging_trade_shares
FROM three_primary 
WHERE exchange_date >= @start_date
AND exchange_date <= @end_date 
AND stock_id = Any(@stock_ids::text[]) 
ORDER BY stock_id, exchange_date desc;

-- name: RetrieveThreePrimaryHistory :many
SELECT stock_id, exchange_date, 
floor(foreign_trade_shares/1000) AS foreign_trade_shares, 
floor(trust_trade_shares/1000) AS trust_trade_shares, 
floor(dealer_trade_shares/1000) AS dealer_trade_shares, 
floor(hedging_trade_shares/1000) AS hedging_trade_shares
FROM three_primary 
WHERE exchange_date >= @start_date
AND exchange_date < @end_date 
AND stock_id = Any(@stock_ids::text[]) 
ORDER BY stock_id, exchange_date desc;
