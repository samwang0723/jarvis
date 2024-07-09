-- name: GetTransactionView :one
SELECT *
FROM transactions
WHERE id = $1;

-- name: UpsertTransactionView :exec
INSERT INTO transactions (id, user_id, order_id, order_type, credit_amount, debit_amount, status, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE
SET user_id = EXCLUDED.user_id, 
  order_id = EXCLUDED.order_id,
  order_type = EXCLUDED.order_type, 
  credit_amount = EXCLUDED.credit_amount, 
  debit_amount = EXCLUDED.debit_amount, 
  status = EXCLUDED.status, 
  version = EXCLUDED.version;
