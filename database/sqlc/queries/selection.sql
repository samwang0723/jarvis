-- name: GetLatestChip :many
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
    floor(d.trade_shares/1000) as volume, 
    floor(COALESCE(t.foreign_trade_shares,0)/1000) as foreignc,
    floor(COALESCE(t.trust_trade_shares,0)/1000) as trust, 
    floor(COALESCE(t.hedging_trade_shares,0)/1000) as hedging,
    floor(COALESCE(t.dealer_trade_shares,0)/1000) as dealer
FROM 
    stake_concentration s
LEFT JOIN 
    stocks c ON c.id = s.stock_id
LEFT JOIN 
    daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = $1)
LEFT JOIN 
    three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = $1)
WHERE 
    (
        CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
    ) >= 4
    AND c.name IS NOT NULL
    AND s.exchange_date = $1
    AND d.trade_shares >= 1000000
ORDER BY 
    s.stock_id;

-- name: GetEligibleStocksFromDate :many
select s.stock_id, c.market
from stake_concentration s
left join stocks c on c.id = s.stock_id
left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = $1)
where (
   CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
) >= 4
AND c.name IS NOT NULL
and s.exchange_date = $1
and d.trade_shares >= 1000000;

-- name: GetEligibleStocksFromPicked :many
select p.stock_id, c.market
from picked_stocks p
left join stocks c on c.id = p.stock_id 
where p.deleted_at is null;

-- name: GetEligibleStocksFromOrder :many
select o.stock_id, c.market
from orders o
left join stocks c on c.id = o.stock_id 
where o.status != 'closed';

-- name: ListSelectionsFromPicked :many
select 
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
  floor(d.trade_shares/1000) as volume, 
  floor(COALESCE(t.foreign_trade_shares,0)/1000) as foreignc,
  floor(COALESCE(t.trust_trade_shares,0)/1000) as trust, 
  floor(COALESCE(t.hedging_trade_shares,0)/1000) as hedging,
  floor(COALESCE(t.dealer_trade_shares,0)/1000) as dealer
from stake_concentration s
left join stocks c on c.id = s.stock_id
left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = $1)
left join three_primary t on (t.stock_id = s.stock_id and t.exchange_date = $1)
where s.stock_id = ANY(@stock_ids::text[])
AND c.name IS NOT NULL
and s.exchange_date = $1
order by s.stock_id;

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
    FLOOR(d.trade_shares/1000) as volume, 
    FLOOR(COALESCE(t.foreign_trade_shares,0)/1000) as foreignc,
    FLOOR(COALESCE(t.trust_trade_shares,0)/1000) as trust, 
    FLOOR(COALESCE(t.hedging_trade_shares,0)/1000) as hedging,
    FLOOR(COALESCE(t.dealer_trade_shares,0)/1000) as dealer,
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
select stock_id, exchange_date, 
floor(foreign_trade_shares/1000) as foreign_trade_shares, 
floor(trust_trade_shares/1000) as trust_trade_shares, 
floor(dealer_trade_shares/1000) as dealer_trade_shares, 
floor(hedging_trade_shares/1000) as hedging_trade_shares
from three_primary where exchange_date >= @start_date
and exchange_date <= @end_date and stock_id = Any(@stock_ids::text[]) 
order by stock_id, exchange_date desc;

-- name: RetrieveThreePrimaryHistory :many
select stock_id, exchange_date, 
floor(foreign_trade_shares/1000) as foreign_trade_shares, 
floor(trust_trade_shares/1000) as trust_trade_shares, 
floor(dealer_trade_shares/1000) as dealer_trade_shares, 
floor(hedging_trade_shares/1000) as hedging_trade_shares
from three_primary where exchange_date >= @start_date
and exchange_date < @end_date and stock_id = Any(@stock_ids::text[]) 
order by stock_id, exchange_date desc;
