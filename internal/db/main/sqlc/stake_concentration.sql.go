// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: stake_concentration.sql

package sqlcdb

import (
	"context"
)

const BatchUpsertStakeConcentration = `-- name: BatchUpsertStakeConcentration :exec
INSERT INTO stake_concentration (
  stock_id, 
  exchange_date, 
  sum_buy_shares, 
  sum_sell_shares, 
  avg_buy_price, 
  avg_sell_price, 
  concentration_1, 
  concentration_5, 
  concentration_10, 
  concentration_20, 
  concentration_60)
VALUES (
  unnest($1::varchar[]), 
  unnest($2::varchar[]), 
  unnest($3::bigint[]), 
  unnest($4::bigint[]), 
  unnest($5::numeric[]),
  unnest($6::numeric[]),
  unnest($7::numeric[]),
  unnest($8::numeric[]),
  unnest($9::numeric[]),
  unnest($10::numeric[]),
  unnest($11::numeric[])
)
ON CONFLICT (stock_id, exchange_date) DO UPDATE
SET sum_buy_shares = EXCLUDED.sum_buy_shares,
    sum_sell_shares = EXCLUDED.sum_sell_shares,
    avg_buy_price = EXCLUDED.avg_buy_price,
    avg_sell_price = EXCLUDED.avg_sell_price,
    concentration_1 = EXCLUDED.concentration_1,
    concentration_5 = EXCLUDED.concentration_5,
    concentration_10 = EXCLUDED.concentration_10,
    concentration_20 = EXCLUDED.concentration_20,
    concentration_60 = EXCLUDED.concentration_60
`

type BatchUpsertStakeConcentrationParams struct {
	StockID         []string
	ExchangeDate    []string
	SumBuyShares    []int64
	SumSellShares   []int64
	AvgBuyPrice     []float64
	AvgSellPrice    []float64
	Concentration1  []float64
	Concentration5  []float64
	Concentration10 []float64
	Concentration20 []float64
	Concentration60 []float64
}

func (q *Queries) BatchUpsertStakeConcentration(ctx context.Context, arg *BatchUpsertStakeConcentrationParams) error {
	_, err := q.db.Exec(ctx, BatchUpsertStakeConcentration,
		arg.StockID,
		arg.ExchangeDate,
		arg.SumBuyShares,
		arg.SumSellShares,
		arg.AvgBuyPrice,
		arg.AvgSellPrice,
		arg.Concentration1,
		arg.Concentration5,
		arg.Concentration10,
		arg.Concentration20,
		arg.Concentration60,
	)
	return err
}

const GetStakeConcentrationByStockID = `-- name: GetStakeConcentrationByStockID :one
SELECT id, stock_id, exchange_date, sum_buy_shares, sum_sell_shares, avg_buy_price, avg_sell_price, concentration_1, concentration_5, concentration_10, concentration_20, concentration_60, created_at, updated_at, deleted_at 
FROM stake_concentration
WHERE stock_id = $1 AND exchange_date = $2
LIMIT 1
`

type GetStakeConcentrationByStockIDParams struct {
	StockID      string
	ExchangeDate string
}

func (q *Queries) GetStakeConcentrationByStockID(ctx context.Context, arg *GetStakeConcentrationByStockIDParams) (*StakeConcentration, error) {
	row := q.db.QueryRow(ctx, GetStakeConcentrationByStockID, arg.StockID, arg.ExchangeDate)
	var i StakeConcentration
	err := row.Scan(
		&i.ID,
		&i.StockID,
		&i.ExchangeDate,
		&i.SumBuyShares,
		&i.SumSellShares,
		&i.AvgBuyPrice,
		&i.AvgSellPrice,
		&i.Concentration1,
		&i.Concentration5,
		&i.Concentration10,
		&i.Concentration20,
		&i.Concentration60,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const GetStakeConcentrationLatestDataPoint = `-- name: GetStakeConcentrationLatestDataPoint :one
SELECT exchange_date FROM stake_concentration
ORDER BY exchange_date DESC LIMIT 1
`

func (q *Queries) GetStakeConcentrationLatestDataPoint(ctx context.Context) (string, error) {
	row := q.db.QueryRow(ctx, GetStakeConcentrationLatestDataPoint)
	var exchange_date string
	err := row.Scan(&exchange_date)
	return exchange_date, err
}

const GetStakeConcentrationsWithVolumes = `-- name: GetStakeConcentrationsWithVolumes :many
SELECT a.trade_shares,
  COALESCE(b.sum_buy_shares, 0)::bigint - COALESCE(b.sum_sell_shares, 0)::bigint AS diff,
  a.exchange_date
FROM daily_closes a
LEFT JOIN stake_concentration b ON (a.stock_id, a.exchange_date) = (b.stock_id, b.exchange_date)
WHERE a.stock_id = $1 
AND a.exchange_date <= $2
ORDER BY a.exchange_date DESC
LIMIT 60
`

type GetStakeConcentrationsWithVolumesParams struct {
	StockID      string
	ExchangeDate string
}

type GetStakeConcentrationsWithVolumesRow struct {
	TradeShares  *int64
	Diff         int32
	ExchangeDate string
}

func (q *Queries) GetStakeConcentrationsWithVolumes(ctx context.Context, arg *GetStakeConcentrationsWithVolumesParams) ([]*GetStakeConcentrationsWithVolumesRow, error) {
	rows, err := q.db.Query(ctx, GetStakeConcentrationsWithVolumes, arg.StockID, arg.ExchangeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetStakeConcentrationsWithVolumesRow
	for rows.Next() {
		var i GetStakeConcentrationsWithVolumesRow
		if err := rows.Scan(&i.TradeShares, &i.Diff, &i.ExchangeDate); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const HasStakeConcentration = `-- name: HasStakeConcentration :one
SELECT EXISTS (
  SELECT 1 FROM stake_concentration
  WHERE exchange_date = $1
)
`

func (q *Queries) HasStakeConcentration(ctx context.Context, exchangeDate string) (bool, error) {
	row := q.db.QueryRow(ctx, HasStakeConcentration, exchangeDate)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
