-- name: GetBalanceView :one
SELECT *
FROM balance_views
WHERE id = $1;

-- name: GetBalanceViewForUpdate :one
SELECT *
FROM balance_views
WHERE id = $1
FOR UPDATE;

-- name: UpsertBalanceView :exec
INSERT INTO balance_views (id, balance, available, pending, version)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
SET balance = EXCLUDED.balance,
    available = EXCLUDED.available,
    pending = EXCLUDED.pending,
    version = EXCLUDED.version;
