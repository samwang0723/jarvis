// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: daily_closes.sql

package sqlcdb

import (
	"context"

	"github.com/ericlagergren/decimal"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const BatchUpsertDailyClose = `-- name: BatchUpsertDailyClose :exec
INSERT INTO daily_closes (
    stock_id, exchange_date, trade_shares, transactions, turnover, open, close, high, low, price_diff
  ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
  )
ON CONFLICT (stock_id, exchange_date) DO UPDATE SET
    trade_shares = EXCLUDED.trade_shares,
    transactions = EXCLUDED.transactions,
    turnover = EXCLUDED.turnover,
    open = EXCLUDED.open,
    close = EXCLUDED.close,
    high = EXCLUDED.high,
    low = EXCLUDED.low,
    price_diff = EXCLUDED.price_diff
`

type BatchUpsertDailyCloseParams struct {
	StockID      string
	ExchangeDate string
	TradeShares  *int64
	Transactions *int64
	Turnover     *int64
	Open         decimal.Big
	Close        decimal.Big
	High         decimal.Big
	Low          decimal.Big
	PriceDiff    decimal.Big
}

func (q *Queries) BatchUpsertDailyClose(ctx context.Context, arg *BatchUpsertDailyCloseParams) error {
	_, err := q.db.Exec(ctx, BatchUpsertDailyClose,
		arg.StockID,
		arg.ExchangeDate,
		arg.TradeShares,
		arg.Transactions,
		arg.Turnover,
		arg.Open,
		arg.Close,
		arg.High,
		arg.Low,
		arg.PriceDiff,
	)
	return err
}

const CreateDailyClose = `-- name: CreateDailyClose :exec
INSERT INTO daily_closes (
    stock_id, exchange_date, trade_shares, transactions, turnover, open, close, high, low, price_diff
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
`

type CreateDailyCloseParams struct {
	StockID      string
	ExchangeDate string
	TradeShares  *int64
	Transactions *int64
	Turnover     *int64
	Open         decimal.Big
	Close        decimal.Big
	High         decimal.Big
	Low          decimal.Big
	PriceDiff    decimal.Big
}

func (q *Queries) CreateDailyClose(ctx context.Context, arg *CreateDailyCloseParams) error {
	_, err := q.db.Exec(ctx, CreateDailyClose,
		arg.StockID,
		arg.ExchangeDate,
		arg.TradeShares,
		arg.Transactions,
		arg.Turnover,
		arg.Open,
		arg.Close,
		arg.High,
		arg.Low,
		arg.PriceDiff,
	)
	return err
}

const HasDailyClose = `-- name: HasDailyClose :one
SELECT stock_id FROM daily_closes
WHERE exchange_date = $1
LIMIT 1
`

func (q *Queries) HasDailyClose(ctx context.Context, exchangeDate string) (string, error) {
	row := q.db.QueryRow(ctx, HasDailyClose, exchangeDate)
	var stock_id string
	err := row.Scan(&stock_id)
	return stock_id, err
}

const ListDailyClose = `-- name: ListDailyClose :many
SELECT t.id, t.stock_id, t.exchange_date, t.transactions,
       FLOOR(t.trade_shares/1000) AS trade_shares, FLOOR(t.turnover/1000) AS turnover,
       t.open, t.high, t.close, t.low, t.price_diff, t.created_at, t.updated_at, t.deleted_at
FROM (
    SELECT daily_closes.id FROM daily_closes
    WHERE daily_closes.exchange_date >= $1 AND daily_closes.stock_id = $2
    ORDER BY daily_closes.exchange_date DESC
    LIMIT $3 OFFSET $4
) q
JOIN daily_closes t ON t.id = q.id
`

type ListDailyCloseParams struct {
	ExchangeDate string
	StockID      string
	Limit        int32
	Offset       int32
}

type ListDailyCloseRow struct {
	ID           uuid.UUID
	StockID      string
	ExchangeDate string
	Transactions *int64
	TradeShares  float64
	Turnover     float64
	Open         decimal.Big
	High         decimal.Big
	Close        decimal.Big
	Low          decimal.Big
	PriceDiff    decimal.Big
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	DeletedAt    pgtype.Timestamp
}

func (q *Queries) ListDailyClose(ctx context.Context, arg *ListDailyCloseParams) ([]*ListDailyCloseRow, error) {
	rows, err := q.db.Query(ctx, ListDailyClose,
		arg.ExchangeDate,
		arg.StockID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListDailyCloseRow
	for rows.Next() {
		var i ListDailyCloseRow
		if err := rows.Scan(
			&i.ID,
			&i.StockID,
			&i.ExchangeDate,
			&i.Transactions,
			&i.TradeShares,
			&i.Turnover,
			&i.Open,
			&i.High,
			&i.Close,
			&i.Low,
			&i.PriceDiff,
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

const ListLatestPrice = `-- name: ListLatestPrice :many
WITH RankedCloses AS (
    SELECT 
        daily_closes.stock_id,
        daily_closes.close,
        daily_closes.exchange_date,
        ROW_NUMBER() OVER (PARTITION BY daily_closes.stock_id ORDER BY daily_closes.exchange_date DESC) AS rn
    FROM 
        daily_closes
    WHERE 
        daily_closes.stock_id = ANY($1::text[])
)
SELECT 
    RankedCloses.stock_id,
    RankedCloses.close,
    RankedCloses.exchange_date
FROM 
    RankedCloses
WHERE 
    rn = 1
`

type ListLatestPriceRow struct {
	StockID      string
	Close        decimal.Big
	ExchangeDate string
}

func (q *Queries) ListLatestPrice(ctx context.Context, stockIds []string) ([]*ListLatestPriceRow, error) {
	rows, err := q.db.Query(ctx, ListLatestPrice, stockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListLatestPriceRow
	for rows.Next() {
		var i ListLatestPriceRow
		if err := rows.Scan(&i.StockID, &i.Close, &i.ExchangeDate); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
