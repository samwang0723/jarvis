// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: selection.sql

package sqlcdb

import (
	"context"
)

const GetEligibleStocksFromDate = `-- name: GetEligibleStocksFromDate :many
select s.stock_id, c.market
from stake_concentration s
left join stocks c on c.id = s.stock_id
left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = $1)
where (
   CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
) >= 4
and s.exchange_date = $1
and d.trade_shares >= 1000000
`

type GetEligibleStocksFromDateRow struct {
	StockID string
	Market  *string
}

func (q *Queries) GetEligibleStocksFromDate(ctx context.Context, exchangeDate string) ([]*GetEligibleStocksFromDateRow, error) {
	rows, err := q.db.Query(ctx, GetEligibleStocksFromDate, exchangeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEligibleStocksFromDateRow
	for rows.Next() {
		var i GetEligibleStocksFromDateRow
		if err := rows.Scan(&i.StockID, &i.Market); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetEligibleStocksFromOrder = `-- name: GetEligibleStocksFromOrder :many
select o.stock_id, c.market
from orders o
left join stocks c on c.id = o.stock_id 
where o.status != 'closed'
`

type GetEligibleStocksFromOrderRow struct {
	StockID string
	Market  *string
}

func (q *Queries) GetEligibleStocksFromOrder(ctx context.Context) ([]*GetEligibleStocksFromOrderRow, error) {
	rows, err := q.db.Query(ctx, GetEligibleStocksFromOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEligibleStocksFromOrderRow
	for rows.Next() {
		var i GetEligibleStocksFromOrderRow
		if err := rows.Scan(&i.StockID, &i.Market); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetEligibleStocksFromPicked = `-- name: GetEligibleStocksFromPicked :many
select p.stock_id, c.market
from picked_stocks p
left join stocks c on c.id = p.stock_id 
where p.deleted_at is null
`

type GetEligibleStocksFromPickedRow struct {
	StockID string
	Market  *string
}

func (q *Queries) GetEligibleStocksFromPicked(ctx context.Context) ([]*GetEligibleStocksFromPickedRow, error) {
	rows, err := q.db.Query(ctx, GetEligibleStocksFromPicked)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEligibleStocksFromPickedRow
	for rows.Next() {
		var i GetEligibleStocksFromPickedRow
		if err := rows.Scan(&i.StockID, &i.Market); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetHighestPrice = `-- name: GetHighestPrice :many
SELECT stock_id, MAX(high)::numeric AS high 
FROM daily_closes 
WHERE exchange_date >= $1 
  AND exchange_date < $2 
  AND stock_id = ANY($3::text[]) 
GROUP BY stock_id
`

type GetHighestPriceParams struct {
	ExchangeDate   string
	ExchangeDate_2 string
	StockIds       []string
}

type GetHighestPriceRow struct {
	StockID string
	High    float64
}

func (q *Queries) GetHighestPrice(ctx context.Context, arg *GetHighestPriceParams) ([]*GetHighestPriceRow, error) {
	rows, err := q.db.Query(ctx, GetHighestPrice, arg.ExchangeDate, arg.ExchangeDate_2, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetHighestPriceRow
	for rows.Next() {
		var i GetHighestPriceRow
		if err := rows.Scan(&i.StockID, &i.High); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetLatestChip = `-- name: GetLatestChip :many
SELECT 
    s.stock_id, 
    c.name, 
  (c.category || '.' || c.market)::text AS category, 
    s.exchange_date, 
    d.open, 
    d.close, 
    d.high, 
    d.low, 
    d.price_diff,
    s.concentration_1, 
    s.concentration_5, 
    s.concentration_10, 
    s.concentration_20, 
    s.concentration_60, 
    floor(d.trade_shares/1000) as volume, 
    floor(t.foreign_trade_shares/1000) as foreignc,
    floor(t.trust_trade_shares/1000) as trust, 
    floor(t.hedging_trade_shares/1000) as hedging,
    floor(t.dealer_trade_shares/1000) as dealer
FROM 
    stake_concentration s
LEFT JOIN 
    stocks c ON c.id = s.stock_id
LEFT JOIN 
    daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = $1)
LEFT JOIN 
    three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = $1)
WHERE 
    (
        CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
        CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
    ) >= 4
    AND s.exchange_date = $1
    AND d.trade_shares >= 1000000
ORDER BY 
    s.stock_id
`

type GetLatestChipRow struct {
	StockID         string
	Name            *string
	Category        string
	ExchangeDate    string
	Open            float64
	Close           float64
	High            float64
	Low             float64
	PriceDiff       float64
	Concentration1  float64
	Concentration5  float64
	Concentration10 float64
	Concentration20 float64
	Concentration60 float64
	Volume          float64
	Foreignc        float64
	Trust           float64
	Hedging         float64
	Dealer          float64
}

func (q *Queries) GetLatestChip(ctx context.Context, exchangeDate string) ([]*GetLatestChipRow, error) {
	rows, err := q.db.Query(ctx, GetLatestChip, exchangeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetLatestChipRow
	for rows.Next() {
		var i GetLatestChipRow
		if err := rows.Scan(
			&i.StockID,
			&i.Name,
			&i.Category,
			&i.ExchangeDate,
			&i.Open,
			&i.Close,
			&i.High,
			&i.Low,
			&i.PriceDiff,
			&i.Concentration1,
			&i.Concentration5,
			&i.Concentration10,
			&i.Concentration20,
			&i.Concentration60,
			&i.Volume,
			&i.Foreignc,
			&i.Trust,
			&i.Hedging,
			&i.Dealer,
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

const GetStartDate = `-- name: GetStartDate :one
SELECT MIN(a.exchange_date)::text 
FROM (
    SELECT exchange_date 
    FROM stake_concentration 
    GROUP BY exchange_date 
    ORDER BY exchange_date DESC 
    LIMIT 120
) AS a
`

func (q *Queries) GetStartDate(ctx context.Context) (string, error) {
	row := q.db.QueryRow(ctx, GetStartDate)
	var column_1 string
	err := row.Scan(&column_1)
	return column_1, err
}

const ListSelections = `-- name: ListSelections :many
WITH average AS (
    SELECT 
        stock_id, 
        AVG(trade_shares) AS avg_volume
    FROM 
        daily_closes
    WHERE 
        exchange_date BETWEEN TO_CHAR(TO_DATE($1, 'YYYYMMDD') - INTERVAL '5' day, 'YYYYMMDD') 
          AND TO_CHAR(TO_DATE($1, 'YYYYMMDD') - INTERVAL '1' day, 'YYYYMMDD')
    GROUP BY 
        stock_id
)
SELECT 
    s.stock_id, 
    c.name, 
    (c.category || '.' || c.market)::text AS category, 
    s.exchange_date, 
    d.open, 
    d.close, 
    d.high, 
    d.low, 
    d.price_diff,
    s.concentration_1, 
    s.concentration_5, 
    s.concentration_10, 
    s.concentration_20, 
    s.concentration_60, 
    FLOOR(d.trade_shares/1000) as volume, 
    FLOOR(t.foreign_trade_shares/1000) as foreignc,
    FLOOR(t.trust_trade_shares/1000) as trust, 
    FLOOR(t.hedging_trade_shares/1000) as hedging,
    FLOOR(t.dealer_trade_shares/1000) as dealer,
    a.avg_volume
FROM 
    stake_concentration s
LEFT JOIN stocks c ON c.id = s.stock_id
LEFT JOIN daily_closes d ON (d.stock_id = s.stock_id AND d.exchange_date = $1)
LEFT JOIN three_primary t ON (t.stock_id = s.stock_id AND t.exchange_date = $1)
LEFT JOIN average a ON a.stock_id = s.stock_id
WHERE (
   CASE WHEN s.concentration_1 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_5 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_10 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_20 > 0 THEN 1 ELSE 0 END +
   CASE WHEN s.concentration_60 > 0 THEN 1 ELSE 0 END
) >= 4
AND s.exchange_date = $1
AND d.trade_shares >= 3000000
AND a.avg_volume >= 1000000
ORDER BY s.stock_id
`

type ListSelectionsRow struct {
	StockID         string
	Name            *string
	Category        string
	ExchangeDate    string
	Open            float64
	Close           float64
	High            float64
	Low             float64
	PriceDiff       float64
	Concentration1  float64
	Concentration5  float64
	Concentration10 float64
	Concentration20 float64
	Concentration60 float64
	Volume          float64
	Foreignc        float64
	Trust           float64
	Hedging         float64
	Dealer          float64
	AvgVolume       *float64
}

func (q *Queries) ListSelections(ctx context.Context, exchangeDate string) ([]*ListSelectionsRow, error) {
	rows, err := q.db.Query(ctx, ListSelections, exchangeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListSelectionsRow
	for rows.Next() {
		var i ListSelectionsRow
		if err := rows.Scan(
			&i.StockID,
			&i.Name,
			&i.Category,
			&i.ExchangeDate,
			&i.Open,
			&i.Close,
			&i.High,
			&i.Low,
			&i.PriceDiff,
			&i.Concentration1,
			&i.Concentration5,
			&i.Concentration10,
			&i.Concentration20,
			&i.Concentration60,
			&i.Volume,
			&i.Foreignc,
			&i.Trust,
			&i.Hedging,
			&i.Dealer,
			&i.AvgVolume,
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

const ListSelectionsFromPicked = `-- name: ListSelectionsFromPicked :many
select s.stock_id, c.name, (c.category || '.' || c.market)::text AS category, s.exchange_date, d.open, 
d.close, d.high, d.low, d.price_diff,s.concentration_1, s.concentration_5, s.concentration_10, 
s.concentration_20, s.concentration_60, floor(d.trade_shares/1000) as volume, 
floor(t.foreign_trade_shares/1000) as foreignc,
floor(t.trust_trade_shares/1000) as trust, floor(t.hedging_trade_shares/1000) as hedging,
floor(t.dealer_trade_shares/1000) as dealer
from stake_concentration s
left join stocks c on c.id = s.stock_id
left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = $1)
left join three_primary t on (t.stock_id = s.stock_id and t.exchange_date = $1)
where s.stock_id = ANY($2::text[])
and s.exchange_date = $1
order by s.stock_id
`

type ListSelectionsFromPickedParams struct {
	ExchangeDate string
	StockIds     []string
}

type ListSelectionsFromPickedRow struct {
	StockID         string
	Name            *string
	Category        string
	ExchangeDate    string
	Open            float64
	Close           float64
	High            float64
	Low             float64
	PriceDiff       float64
	Concentration1  float64
	Concentration5  float64
	Concentration10 float64
	Concentration20 float64
	Concentration60 float64
	Volume          float64
	Foreignc        float64
	Trust           float64
	Hedging         float64
	Dealer          float64
}

func (q *Queries) ListSelectionsFromPicked(ctx context.Context, arg *ListSelectionsFromPickedParams) ([]*ListSelectionsFromPickedRow, error) {
	rows, err := q.db.Query(ctx, ListSelectionsFromPicked, arg.ExchangeDate, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListSelectionsFromPickedRow
	for rows.Next() {
		var i ListSelectionsFromPickedRow
		if err := rows.Scan(
			&i.StockID,
			&i.Name,
			&i.Category,
			&i.ExchangeDate,
			&i.Open,
			&i.Close,
			&i.High,
			&i.Low,
			&i.PriceDiff,
			&i.Concentration1,
			&i.Concentration5,
			&i.Concentration10,
			&i.Concentration20,
			&i.Concentration60,
			&i.Volume,
			&i.Foreignc,
			&i.Trust,
			&i.Hedging,
			&i.Dealer,
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

const RetrieveDailyCloseHistory = `-- name: RetrieveDailyCloseHistory :many
SELECT stock_id, exchange_date, close, trade_shares 
FROM daily_closes 
WHERE exchange_date >= $1 
  AND exchange_date < $2 
  AND stock_id = ANY($3::text[]) 
ORDER BY stock_id, exchange_date DESC
`

type RetrieveDailyCloseHistoryParams struct {
	ExchangeDate   string
	ExchangeDate_2 string
	StockIds       []string
}

type RetrieveDailyCloseHistoryRow struct {
	StockID      string
	ExchangeDate string
	Close        float64
	TradeShares  *int64
}

func (q *Queries) RetrieveDailyCloseHistory(ctx context.Context, arg *RetrieveDailyCloseHistoryParams) ([]*RetrieveDailyCloseHistoryRow, error) {
	rows, err := q.db.Query(ctx, RetrieveDailyCloseHistory, arg.ExchangeDate, arg.ExchangeDate_2, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*RetrieveDailyCloseHistoryRow
	for rows.Next() {
		var i RetrieveDailyCloseHistoryRow
		if err := rows.Scan(
			&i.StockID,
			&i.ExchangeDate,
			&i.Close,
			&i.TradeShares,
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

const RetrieveDailyCloseHistoryWithDate = `-- name: RetrieveDailyCloseHistoryWithDate :many
SELECT stock_id, exchange_date, close, trade_shares 
FROM daily_closes 
WHERE exchange_date >= $1 
  AND exchange_date <= $2 
  AND stock_id = ANY($3::text[]) 
ORDER BY stock_id, exchange_date DESC
`

type RetrieveDailyCloseHistoryWithDateParams struct {
	ExchangeDate   string
	ExchangeDate_2 string
	StockIds       []string
}

type RetrieveDailyCloseHistoryWithDateRow struct {
	StockID      string
	ExchangeDate string
	Close        float64
	TradeShares  *int64
}

func (q *Queries) RetrieveDailyCloseHistoryWithDate(ctx context.Context, arg *RetrieveDailyCloseHistoryWithDateParams) ([]*RetrieveDailyCloseHistoryWithDateRow, error) {
	rows, err := q.db.Query(ctx, RetrieveDailyCloseHistoryWithDate, arg.ExchangeDate, arg.ExchangeDate_2, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*RetrieveDailyCloseHistoryWithDateRow
	for rows.Next() {
		var i RetrieveDailyCloseHistoryWithDateRow
		if err := rows.Scan(
			&i.StockID,
			&i.ExchangeDate,
			&i.Close,
			&i.TradeShares,
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

const RetrieveThreePrimaryHistory = `-- name: RetrieveThreePrimaryHistory :many
select stock_id, exchange_date, 
floor(foreign_trade_shares/1000) as foreign_trade_shares, 
floor(trust_trade_shares/1000) as trust_trade_shares, 
floor(dealer_trade_shares/1000) as dealer_trade_shares, 
floor(hedging_trade_shares/1000) as hedging_trade_shares
from three_primary where exchange_date >= $1
and exchange_date < $2 and stock_id = Any($3::text[]) 
order by stock_id, exchange_date desc
`

type RetrieveThreePrimaryHistoryParams struct {
	ExchangeDate   string
	ExchangeDate_2 string
	StockIds       []string
}

type RetrieveThreePrimaryHistoryRow struct {
	StockID            string
	ExchangeDate       string
	ForeignTradeShares float64
	TrustTradeShares   float64
	DealerTradeShares  float64
	HedgingTradeShares float64
}

func (q *Queries) RetrieveThreePrimaryHistory(ctx context.Context, arg *RetrieveThreePrimaryHistoryParams) ([]*RetrieveThreePrimaryHistoryRow, error) {
	rows, err := q.db.Query(ctx, RetrieveThreePrimaryHistory, arg.ExchangeDate, arg.ExchangeDate_2, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*RetrieveThreePrimaryHistoryRow
	for rows.Next() {
		var i RetrieveThreePrimaryHistoryRow
		if err := rows.Scan(
			&i.StockID,
			&i.ExchangeDate,
			&i.ForeignTradeShares,
			&i.TrustTradeShares,
			&i.DealerTradeShares,
			&i.HedgingTradeShares,
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

const RetrieveThreePrimaryHistoryWithDate = `-- name: RetrieveThreePrimaryHistoryWithDate :many
select stock_id, exchange_date, 
floor(foreign_trade_shares/1000) as foreign_trade_shares, 
floor(trust_trade_shares/1000) as trust_trade_shares, 
floor(dealer_trade_shares/1000) as dealer_trade_shares, 
floor(hedging_trade_shares/1000) as hedging_trade_shares
from three_primary where exchange_date >= $1
and exchange_date <= $2 and stock_id = Any($3::text[]) 
order by stock_id, exchange_date desc
`

type RetrieveThreePrimaryHistoryWithDateParams struct {
	ExchangeDate   string
	ExchangeDate_2 string
	StockIds       []string
}

type RetrieveThreePrimaryHistoryWithDateRow struct {
	StockID            string
	ExchangeDate       string
	ForeignTradeShares float64
	TrustTradeShares   float64
	DealerTradeShares  float64
	HedgingTradeShares float64
}

func (q *Queries) RetrieveThreePrimaryHistoryWithDate(ctx context.Context, arg *RetrieveThreePrimaryHistoryWithDateParams) ([]*RetrieveThreePrimaryHistoryWithDateRow, error) {
	rows, err := q.db.Query(ctx, RetrieveThreePrimaryHistoryWithDate, arg.ExchangeDate, arg.ExchangeDate_2, arg.StockIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*RetrieveThreePrimaryHistoryWithDateRow
	for rows.Next() {
		var i RetrieveThreePrimaryHistoryWithDateRow
		if err := rows.Scan(
			&i.StockID,
			&i.ExchangeDate,
			&i.ForeignTradeShares,
			&i.TrustTradeShares,
			&i.DealerTradeShares,
			&i.HedgingTradeShares,
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
