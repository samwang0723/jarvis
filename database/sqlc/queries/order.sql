-- name: GetOrder :one
SELECT *
FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE user_id = $1
  AND (stock_id = ANY(@stock_ids::text[]) OR NOT @filter_by_stock_id::bool)
  AND (@status::VARCHAR = '' OR status = @status)
  AND (@exchange_month::VARCHAR = '' 
    OR sell_exchange_date LIKE @exchange_month::VARCHAR || '%' 
    OR buy_exchange_date LIKE @exchange_month::VARCHAR || '%')
LIMIT $2 OFFSET $3;

-- name: ListOpenOrders :many
SELECT id 
FROM orders 
WHERE user_id = $1 
  AND stock_id = $2 
  AND status IN ('created', 'changed') 
  AND (
    (@order_type::VARCHAR = 'Sell' AND buy_quantity - sell_quantity > 0) 
    OR 
    (@order_type::VARCHAR = 'Buy' AND sell_quantity - buy_quantity > 0)
  )
ORDER BY created_at ASC;


-- name: UpsertOrder :exec
INSERT INTO orders (id, user_id, stock_id, buy_price, buy_quantity,
buy_exchange_date, sell_price, sell_quantity, sell_exchange_date, profitable_price,
status, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (id) DO UPDATE
SET user_id = EXCLUDED.user_id, 
  stock_id = EXCLUDED.stock_id,
  buy_price = EXCLUDED.buy_price, 
  buy_quantity = EXCLUDED.buy_quantity, 
  buy_exchange_date = EXCLUDED.buy_exchange_date,
  sell_price = EXCLUDED.sell_price, 
  sell_quantity = EXCLUDED.sell_quantity, 
  sell_exchange_date = EXCLUDED.sell_exchange_date,
  profitable_price = EXCLUDED.profitable_price, 
  status = EXCLUDED.status, 
  version = EXCLUDED.version;
