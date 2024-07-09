-- name: CreatePickedStocks :exec
INSERT INTO picked_stocks (user_id, stock_id)
SELECT unnest(@user_ids::uuid[]), unnest(@stock_ids::text[]);

-- name: DeletePickedStock :exec
UPDATE picked_stocks SET deleted_at = NOW() 
WHERE user_id = $1 AND stock_id = $2 AND deleted_at IS NULL;

-- name: ListPickedStocks :many
SELECT * FROM picked_stocks 
WHERE deleted_at IS NULL AND user_id = $1;
