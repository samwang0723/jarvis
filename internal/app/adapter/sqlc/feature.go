package sqlc

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) ListCategories(ctx context.Context) ([]*string, error) {
	return repo.primary().ListCategories(ctx)
}

func (repo *Repo) ListStocks(
	ctx context.Context,
	arg *domain.ListStocksParams,
) ([]*domain.Stock, error) {
	res, err := repo.primary().ListStocks(ctx, toSqlcListStocksParams(arg))
	if err != nil {
		return nil, err
	}
	return fromSqlcStocks(res), nil
}

func toSqlcListStocksParams(arg *domain.ListStocksParams) *sqlcdb.ListStocksParams {
	return &sqlcdb.ListStocksParams{
		Limit:           arg.Limit,
		Offset:          arg.Offset,
		Country:         arg.Country,
		StockIds:        arg.StockIDs,
		FilterByStockID: arg.FilterByStockID,
		Name:            arg.Name,
		Category:        arg.Category,
	}
}

func fromSqlcStocks(stocks []*sqlcdb.Stock) []*domain.Stock {
	result := make([]*domain.Stock, 0, len(stocks))
	for _, stock := range stocks {
		result = append(result, fromSqlcStock(stock))
	}
	return result
}

func fromSqlcStock(stock *sqlcdb.Stock) *domain.Stock {
	return &domain.Stock{
		ID:       stock.ID,
		Name:     stock.Name,
		Country:  stock.Country,
		Category: *stock.Category,
		Market:   *stock.Market,
		Time: domain.Time{
			CreatedAt: &stock.CreatedAt.Time,
			UpdatedAt: &stock.UpdatedAt.Time,
			DeletedAt: &stock.DeletedAt.Time,
		},
	}
}

func (repo *Repo) ListThreePrimary(
	ctx context.Context,
	arg *domain.ListThreePrimaryParams,
) ([]*domain.ThreePrimary, error) {
	res, err := repo.primary().ListThreePrimary(ctx, toSqlcListThreePrimaryParams(arg))
	if err != nil {
		return nil, err
	}
	return fromSqlcThreePrimarys(res), nil
}

func toSqlcListThreePrimaryParams(
	arg *domain.ListThreePrimaryParams,
) *sqlcdb.ListThreePrimaryParams {
	return &sqlcdb.ListThreePrimaryParams{
		Limit:     arg.Limit,
		Offset:    arg.Offset,
		StockID:   arg.StockID,
		StartDate: arg.StartDate,
		EndDate:   arg.EndDate,
	}
}

func fromSqlcThreePrimarys(threePrimary []*sqlcdb.ThreePrimary) []*domain.ThreePrimary {
	result := make([]*domain.ThreePrimary, 0, len(threePrimary))
	for _, tp := range threePrimary {
		result = append(result, fromSqlcThreePrimary(tp))
	}
	return result
}

func fromSqlcThreePrimary(tp *sqlcdb.ThreePrimary) *domain.ThreePrimary {
	return &domain.ThreePrimary{
		ID: domain.ID{
			ID: tp.ID,
		},
		StockID:            tp.StockID,
		ExchangeDate:       tp.ExchangeDate,
		ForeignTradeShares: tp.ForeignTradeShares,
		TrustTradeShares:   tp.TrustTradeShares,
		DealerTradeShares:  tp.DealerTradeShares,
		HedgingTradeShares: tp.HedgingTradeShares,
		Time: domain.Time{
			CreatedAt: &tp.CreatedAt.Time,
			UpdatedAt: &tp.UpdatedAt.Time,
			DeletedAt: &tp.DeletedAt.Time,
		},
	}
}
