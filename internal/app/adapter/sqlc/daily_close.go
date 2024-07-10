package sqlc

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
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
		ExchangeDate: obj.Date,
		TradeShares:  &obj.TradedShares,
		Transactions: &obj.Transactions,
		Turnover:     &obj.Turnover,
		Open:         float64(obj.Open),
		Close:        float64(obj.Close),
		High:         float64(obj.High),
		Low:          float64(obj.Low),
		PriceDiff:    float64(obj.PriceDiff),
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
		Open:         make([]float64, 0, len(dailyClose)),
		Close:        make([]float64, 0, len(dailyClose)),
		High:         make([]float64, 0, len(dailyClose)),
		Low:          make([]float64, 0, len(dailyClose)),
		PriceDiff:    make([]float64, 0, len(dailyClose)),
	}
	for _, dc := range dailyClose {
		result.StockID = append(result.StockID, dc.StockID)
		result.ExchangeDate = append(result.ExchangeDate, dc.Date)
		result.TradeShares = append(result.TradeShares, dc.TradedShares)
		result.Transactions = append(result.Transactions, dc.Transactions)
		result.Turnover = append(result.Turnover, dc.Turnover)
		result.Open = append(result.Open, float64(dc.Open))
		result.Close = append(result.Close, float64(dc.Close))
		result.High = append(result.High, float64(dc.High))
		result.Low = append(result.Low, float64(dc.Low))
		result.PriceDiff = append(result.PriceDiff, float64(dc.PriceDiff))
	}

	return result
}

func toDomainDailyCloseList(res []*sqlcdb.ListDailyCloseRow) []*domain.DailyClose {
	result := make([]*domain.DailyClose, 0, len(res))
	for _, r := range res {
		time := domain.Time{
			CreatedAt: &r.CreatedAt.Time,
			UpdatedAt: &r.UpdatedAt.Time,
		}
		if r.DeletedAt.Valid {
			time.DeletedAt = &r.DeletedAt.Time
		}
		result = append(result, &domain.DailyClose{
			ID: domain.ID{
				ID: r.ID,
			},
			StockID:      r.StockID,
			Date:         r.ExchangeDate,
			TradedShares: int64(r.TradeShares),
			Transactions: *r.Transactions,
			Turnover:     int64(r.Turnover),
			Open:         float32(r.Open),
			Close:        float32(r.Close),
			High:         float32(r.High),
			Low:          float32(r.Low),
			PriceDiff:    float32(r.PriceDiff),
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
			Price:        float32(r.Close),
		})
	}
	return result
}
