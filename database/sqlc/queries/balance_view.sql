-- name: GetBalanceView :one
SELECT id, 
  regexp_replace(balance::text, '[^\d.-]', '', 'g')::numeric as balance, 
  regexp_replace(available::text, '[^\d.-]', '', 'g')::numeric as available, 
  regexp_replace(pending::text, '[^\d.-]', '', 'g')::numeric as pending, 
  version, created_at, updated_at
FROM balance_views
WHERE id = $1;

-- name: UpsertBalanceView :exec
INSERT INTO balance_views (id, balance, available, pending, version)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
SET balance = EXCLUDED.balance,
    available = EXCLUDED.available,
    pending = EXCLUDED.pending,
    version = EXCLUDED.version;
