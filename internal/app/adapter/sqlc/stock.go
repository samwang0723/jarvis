package sqlc

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) BatchUpsertStocks(
	ctx context.Context,
	objs []*domain.Stock,
) error {
	return repo.primary().BatchUpsertStocks(ctx, toSqlcBatchUpsertStocksParams(objs))
}

func (repo *Repo) CreateStock(
	ctx context.Context,
	obj *domain.Stock,
) error {
	return repo.primary().CreateStock(ctx, toSqlcCreateStockParams(obj))
}

func (repo *Repo) DeleteStockByID(
	ctx context.Context,
	id string,
) error {
	return repo.primary().DeleteStockByID(ctx, id)
}

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

func toSqlcBatchUpsertStocksParams(stocks []*domain.Stock) *sqlcdb.BatchUpsertStocksParams {
	result := &sqlcdb.BatchUpsertStocksParams{
		ID:       make([]string, 0, len(stocks)),
		Name:     make([]string, 0, len(stocks)),
		Country:  make([]string, 0, len(stocks)),
		Category: make([]string, 0, len(stocks)),
		Market:   make([]string, 0, len(stocks)),
	}
	for _, stock := range stocks {
		result.ID = append(result.ID, stock.ID)
		result.Name = append(result.Name, stock.Name)
		result.Country = append(result.Country, stock.Country)
		result.Category = append(result.Category, stock.Category)
		result.Market = append(result.Market, stock.Market)
	}
	return result
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

func toSqlcCreateStockParams(stock *domain.Stock) *sqlcdb.CreateStockParams {
	return &sqlcdb.CreateStockParams{
		ID:       stock.ID,
		Name:     stock.Name,
		Country:  stock.Country,
		Category: &stock.Category,
		Market:   &stock.Market,
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
