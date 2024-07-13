package sqlc

import (
	"context"
	"database/sql"

	"github.com/ericlagergren/decimal"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/helper"
)

func (repo *Repo) BatchUpsertDailyClose(
	ctx context.Context,
	objs []*domain.DailyClose,
) error {
	return repo.primary().BatchUpsertDailyClose(ctx, toSqlcBatchUpsertDailyCloseParams(objs))
}

func (repo *Repo) CreateDailyClose(
	ctx context.Context,
	obj *domain.DailyClose,
) error {
	return repo.primary().CreateDailyClose(ctx, &sqlcdb.CreateDailyCloseParams{
		StockID:      obj.StockID,
		ExchangeDate: obj.ExchangeDate,
		TradeShares:  sql.NullInt64{Int64: obj.TradedShares, Valid: true},
		Transactions: sql.NullInt64{Int64: obj.Transactions, Valid: true},
		Turnover:     sql.NullInt64{Int64: obj.Turnover, Valid: true},
		Open:         helper.Float32ToDecimal(obj.Open),
		Close:        helper.Float32ToDecimal(obj.Close),
		High:         helper.Float32ToDecimal(obj.High),
		Low:          helper.Float32ToDecimal(obj.Low),
		PriceDiff:    helper.Float32ToDecimal(obj.PriceDiff),
	})
}

func (repo *Repo) HasDailyClose(ctx context.Context, date string) (bool, error) {
	return repo.primary().HasDailyClose(ctx, date)
}

func (repo *Repo) ListDailyClose(
	ctx context.Context,
	arg *domain.ListDailyCloseParams,
) ([]*domain.DailyClose, error) {
	result, err := repo.primary().ListDailyClose(ctx, &sqlcdb.ListDailyCloseParams{
		Limit:     arg.Limit,
		Offset:    arg.Offset,
		StartDate: arg.StartDate,
		StockID:   arg.StockID,
		EndDate:   arg.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return toDomainDailyCloseList(result), nil
}

func (repo *Repo) ListLatestPrice(
	ctx context.Context,
	stockIDs []string,
) ([]*domain.StockPrice, error) {
	result, err := repo.primary().ListLatestPrice(ctx, stockIDs)
	if err != nil {
		return nil, err
	}
	return toDomainStockPriceList(result), nil
}

func toSqlcBatchUpsertDailyCloseParams(
	dailyClose []*domain.DailyClose,
) *sqlcdb.BatchUpsertDailyCloseParams {
	result := &sqlcdb.BatchUpsertDailyCloseParams{
		StockID:      make([]string, 0, len(dailyClose)),
		ExchangeDate: make([]string, 0, len(dailyClose)),
		TradeShares:  make([]int64, 0, len(dailyClose)),
		Transactions: make([]int64, 0, len(dailyClose)),
		Turnover:     make([]int64, 0, len(dailyClose)),
		Open:         make([]decimal.Big, 0, len(dailyClose)),
		Close:        make([]decimal.Big, 0, len(dailyClose)),
		High:         make([]decimal.Big, 0, len(dailyClose)),
		Low:          make([]decimal.Big, 0, len(dailyClose)),
		PriceDiff:    make([]decimal.Big, 0, len(dailyClose)),
	}
	for _, dc := range dailyClose {
		result.StockID = append(result.StockID, dc.StockID)
		result.ExchangeDate = append(result.ExchangeDate, dc.ExchangeDate)
		result.TradeShares = append(result.TradeShares, dc.TradedShares)
		result.Transactions = append(result.Transactions, dc.Transactions)
		result.Turnover = append(result.Turnover, dc.Turnover)
		result.Open = append(result.Open, helper.Float32ToDecimal(dc.Open))
		result.Close = append(result.Close, helper.Float32ToDecimal(dc.Close))
		result.High = append(result.High, helper.Float32ToDecimal(dc.High))
		result.Low = append(result.Low, helper.Float32ToDecimal(dc.Low))
		result.PriceDiff = append(result.PriceDiff, helper.Float32ToDecimal(dc.PriceDiff))
	}

	return result
}

func toDomainDailyCloseList(res []*sqlcdb.ListDailyCloseRow) []*domain.DailyClose {
	result := make([]*domain.DailyClose, 0, len(res))
	for _, r := range res {
		time := domain.Time{
			CreatedAt: &r.CreatedAt,
			UpdatedAt: &r.UpdatedAt,
		}
		if r.DeletedAt.Valid {
			time.DeletedAt = &r.DeletedAt.Time
		}
		result = append(result, &domain.DailyClose{
			ID: domain.ID{
				ID: r.ID,
			},
			StockID:      r.StockID,
			ExchangeDate: r.ExchangeDate,
			TradedShares: int64(r.TradeShares),
			Transactions: r.Transactions.Int64,
			Turnover:     int64(r.Turnover),
			Open:         helper.DecimalToFloat32(r.Open),
			Close:        helper.DecimalToFloat32(r.Close),
			High:         helper.DecimalToFloat32(r.High),
			Low:          helper.DecimalToFloat32(r.Low),
			PriceDiff:    helper.DecimalToFloat32(r.PriceDiff),
			Time:         time,
		})
	}
	return result
}

func toDomainStockPriceList(res []*sqlcdb.ListLatestPriceRow) []*domain.StockPrice {
	result := make([]*domain.StockPrice, 0, len(res))
	for _, r := range res {
		result = append(result, &domain.StockPrice{
			ExchangeDate: r.ExchangeDate,
			StockID:      r.StockID,
			Price:        helper.DecimalToFloat32(r.Close),
		})
	}
	return result
}
