// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: three_primary.sql

package sqlcdb

import (
	"context"
)

const BatchUpsertThreePrimary = `-- name: BatchUpsertThreePrimary :exec
INSERT INTO three_primary (
    stock_id, exchange_date, foreign_trade_shares, trust_trade_shares, dealer_trade_shares, hedging_trade_shares
) VALUES (
    unnest($1::varchar[]), 
    unnest($2::varchar[]), 
    unnest($3::bigint[]), 
    unnest($4::bigint[]), 
    unnest($5::bigint[]),
    unnest($6::bigint[])
)
ON CONFLICT (stock_id, exchange_date) DO UPDATE
SET
    foreign_trade_shares = EXCLUDED.foreign_trade_shares,
    trust_trade_shares = EXCLUDED.trust_trade_shares,
    dealer_trade_shares = EXCLUDED.dealer_trade_shares,
    hedging_trade_shares = EXCLUDED.hedging_trade_shares
`

type BatchUpsertThreePrimaryParams struct {
	StockID            []string
	ExchangeDate       []string
	ForeignTradeShares []int64
	TrustTradeShares   []int64
	DealerTradeShares  []int64
	HedgingTradeShares []int64
}

func (q *Queries) BatchUpsertThreePrimary(ctx context.Context, arg *BatchUpsertThreePrimaryParams) error {
	_, err := q.db.Exec(ctx, BatchUpsertThreePrimary,
		arg.StockID,
		arg.ExchangeDate,
		arg.ForeignTradeShares,
		arg.TrustTradeShares,
		arg.DealerTradeShares,
		arg.HedgingTradeShares,
	)
	return err
}

const CreateThreePrimary = `-- name: CreateThreePrimary :exec
INSERT INTO three_primary (
    stock_id, exchange_date, foreign_trade_shares, trust_trade_shares, dealer_trade_shares, hedging_trade_shares
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type CreateThreePrimaryParams struct {
	StockID            string
	ExchangeDate       string
	ForeignTradeShares *int64
	TrustTradeShares   *int64
	DealerTradeShares  *int64
	HedgingTradeShares *int64
}

func (q *Queries) CreateThreePrimary(ctx context.Context, arg *CreateThreePrimaryParams) error {
	_, err := q.db.Exec(ctx, CreateThreePrimary,
		arg.StockID,
		arg.ExchangeDate,
		arg.ForeignTradeShares,
		arg.TrustTradeShares,
		arg.DealerTradeShares,
		arg.HedgingTradeShares,
	)
	return err
}

const ListThreePrimary = `-- name: ListThreePrimary :many
WITH filtered AS (
    SELECT three_primary.id
    FROM three_primary
    WHERE three_primary.stock_id = $3 AND three_primary.exchange_date >= $4
    AND ($5::text = '' OR three_primary.exchange_date <= $5::text)
    ORDER BY three_primary.exchange_date DESC
    LIMIT $1 OFFSET $2
)
SELECT t.id, t.stock_id, t.exchange_date, t.foreign_trade_shares, t.trust_trade_shares, t.dealer_trade_shares, t.hedging_trade_shares, t.created_at, t.updated_at, t.deleted_at
FROM filtered f
JOIN three_primary t ON t.id = f.id
`

type ListThreePrimaryParams struct {
	Limit     int32
	Offset    int32
	StockID   string
	StartDate string
	EndDate   string
}

func (q *Queries) ListThreePrimary(ctx context.Context, arg *ListThreePrimaryParams) ([]*ThreePrimary, error) {
	rows, err := q.db.Query(ctx, ListThreePrimary,
		arg.Limit,
		arg.Offset,
		arg.StockID,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ThreePrimary
	for rows.Next() {
		var i ThreePrimary
		if err := rows.Scan(
			&i.ID,
			&i.StockID,
			&i.ExchangeDate,
			&i.ForeignTradeShares,
			&i.TrustTradeShares,
			&i.DealerTradeShares,
			&i.HedgingTradeShares,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
