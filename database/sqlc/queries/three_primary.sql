-- name: CreateThreePrimary :exec
INSERT INTO three_primary (
    stock_id, exchange_date, foreign_trade_shares, trust_trade_shares, dealer_trade_shares, hedging_trade_shares
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: BatchUpsertThreePrimary :exec
INSERT INTO three_primary (
    stock_id, exchange_date, foreign_trade_shares, trust_trade_shares, dealer_trade_shares, hedging_trade_shares
) VALUES (
    unnest($1::varchar[]), unnest($2::varchar[]), unnest($3::bigint[]), unnest($4::bigint[]), unnest($5::bigint[]), unnest($6::bigint[])
)
ON CONFLICT (stock_id, exchange_date) DO UPDATE
SET
    foreign_trade_shares = EXCLUDED.foreign_trade_shares,
    trust_trade_shares = EXCLUDED.trust_trade_shares,
    dealer_trade_shares = EXCLUDED.dealer_trade_shares,
    hedging_trade_shares = EXCLUDED.hedging_trade_shares;

-- name: ListThreePrimary :many
WITH filtered AS (
    SELECT three_primary.id
    FROM three_primary
    WHERE three_primary.stock_id = @stock_id AND three_primary.exchange_date >= @start_date
    AND (@end_date::text = '' OR three_primary.exchange_date <= @end_date::text)
    ORDER BY three_primary.exchange_date DESC
    LIMIT $1 OFFSET $2
)
SELECT t.*
FROM filtered f
JOIN three_primary t ON t.id = f.id;

-- name: CountThreePrimary :one
SELECT COUNT(three_primary.*)
FROM three_primary
WHERE three_primary.stock_id = @stock_id AND three_primary.exchange_date >= @start_date
AND (@end_date::text = '' OR three_primary.exchange_date <= @end_date::text);

