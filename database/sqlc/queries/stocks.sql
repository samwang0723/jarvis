-- name: CreateStock :exec
INSERT INTO stocks (id, stock_id, name, country, category, market)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateStock :exec
UPDATE stocks
SET 
    name = $2,
    country = $3,
    category = $4,
    market = $5
WHERE stock_id = $1;

-- name: BatchUpsertStocks :exec
INSERT INTO stocks (id, stock_id, name, country, category, market)
VALUES (
    unnest($1::int8[]), 
    unnest($2::text[]), 
    unnest($3::text[]), 
    unnest($4::text[]), 
    unnest($5::text[]),
    unnest($6::text[])
)
ON CONFLICT (stock_id) DO UPDATE 
SET 
    name = EXCLUDED.name,
    country = EXCLUDED.country,
    category = EXCLUDED.category,
    market = EXCLUDED.market;

-- name: DeleteStockbyID :exec
UPDATE stocks SET deleted_at = NOW() WHERE id = $1;

-- name: DeleteStockbyStockID :one
SELECT id, stock_id, name, country, category, market, created_at, updated_at, deleted_at
FROM stocks
WHERE stock_id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CountStocks :one
SELECT COUNT(*) FROM stocks
WHERE
    (@country::VARCHAR = '' OR country = @country)
    AND (stock_id = ANY(@stock_ids::text[]) OR NOT @filter_by_stock_id::bool)
    AND (@name::VARCHAR = '' OR name ILIKE '%' || @name || '%')
    AND (@category::VARCHAR = '' OR category = @category);

-- name: ListStocks :many
SELECT * FROM stocks
WHERE
    (@country::VARCHAR = '' OR country = @country)
    AND (stock_id = ANY(@stock_ids::text[]) OR NOT @filter_by_stock_id::bool)
    AND (@name::VARCHAR = '' OR name ILIKE '%' || @name || '%')
    AND (@category::VARCHAR = '' OR category = @category)
ORDER BY stock_id
LIMIT $1 OFFSET $2;

-- name: ListCategories :many
SELECT DISTINCT category FROM stocks;
