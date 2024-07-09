-- name: GetOrderView :one
SELECT *
FROM orders
WHERE id = $1;

-- name: UpsertOrderView :exec
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
