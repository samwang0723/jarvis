-- name: CreateStock :exec
INSERT INTO stocks (id, name, country, category, market)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateStock :exec
UPDATE stocks
SET 
    name = $2,
    country = $3,
    category = $4,
    market = $5
WHERE id = $1;

-- name: BatchUpsertStocks :exec
INSERT INTO stocks (id, name, country, category, market)
VALUES (
    unnest(@id::text[]), 
    unnest(@name::text[]), 
    unnest(@country::text[]), 
    unnest(@category::text[]),
    unnest(@market::text[])
)
ON CONFLICT (id) DO UPDATE 
SET 
    name = EXCLUDED.name,
    country = EXCLUDED.country,
    category = EXCLUDED.category,
    market = EXCLUDED.market;

-- name: DeleteStockByID :exec
UPDATE stocks SET deleted_at = NOW() WHERE id = $1;

-- name: ListStocks :many
SELECT id, name, country, category, market, created_at, updated_at, deleted_at FROM stocks
WHERE
    (@country::VARCHAR = '' OR country = @country)
    AND (id = ANY(@stock_ids::text[]) OR NOT @filter_by_stock_id::bool)
    AND (@name::VARCHAR = '' OR name ILIKE '%' || @name || '%')
    AND (@category::VARCHAR = '' OR category = @category)
    AND deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListCategories :many
SELECT DISTINCT category FROM stocks;
