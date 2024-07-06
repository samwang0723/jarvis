package adapter

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/adapter/sqlc"
	"github.com/samwang0723/jarvis/internal/app/domain"
)

type Adapter interface {
	BatchUpsertStocks(ctx context.Context, objs []*domain.Stock) error
	CreateStock(ctx context.Context, obj *domain.Stock) error
	DeleteStockByID(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]*string, error)
	ListStocks(ctx context.Context, arg *domain.ListStocksParams) ([]*domain.Stock, error)
	BatchUpsertThreePrimary(ctx context.Context, objs []*domain.ThreePrimary) error
	CreateThreePrimary(ctx context.Context, arg *domain.ThreePrimary) error
	ListThreePrimary(
		ctx context.Context,
		arg *domain.ListThreePrimaryParams,
	) ([]*domain.ThreePrimary, error)
	BatchUpsertDailyClose(ctx context.Context, objs []*domain.DailyClose) error
	CreateDailyClose(ctx context.Context, obj *domain.DailyClose) error
	HasDailyClose(ctx context.Context, date string) (bool, error)
	ListDailyClose(
		ctx context.Context,
		arg *domain.ListDailyCloseParams,
	) ([]*domain.DailyClose, error)
	ListLatestPrice(ctx context.Context, stockIDs []string) ([]*domain.StockPrice, error)
}

var _ Adapter = (*Imp)(nil)

type Imp struct {
	repo *sqlc.Repo
}

func NewAdapterImp(repo *sqlc.Repo) *Imp {
	return &Imp{
		repo: repo,
	}
}

func (a *Imp) BatchUpsertStocks(ctx context.Context, objs []*domain.Stock) error {
	return a.repo.BatchUpsertStocks(ctx, objs)
}

func (a *Imp) CreateStock(ctx context.Context, obj *domain.Stock) error {
	return a.repo.CreateStock(ctx, obj)
}

func (a *Imp) DeleteStockByID(ctx context.Context, id string) error {
	return a.repo.DeleteStockByID(ctx, id)
}

func (a *Imp) ListCategories(ctx context.Context) ([]*string, error) {
	return a.repo.ListCategories(ctx)
}

func (a *Imp) ListStocks(
	ctx context.Context,
	arg *domain.ListStocksParams,
) ([]*domain.Stock, error) {
	return a.repo.ListStocks(ctx, arg)
}

func (a *Imp) BatchUpsertThreePrimary(ctx context.Context, objs []*domain.ThreePrimary) error {
	return a.repo.BatchUpsertThreePrimary(ctx, objs)
}

func (a *Imp) CreateThreePrimary(ctx context.Context, arg *domain.ThreePrimary) error {
	return a.repo.CreateThreePrimary(ctx, arg)
}

func (a *Imp) ListThreePrimary(
	ctx context.Context,
	arg *domain.ListThreePrimaryParams,
) ([]*domain.ThreePrimary, error) {
	return a.repo.ListThreePrimary(ctx, arg)
}

func (a *Imp) BatchUpsertDailyClose(ctx context.Context, objs []*domain.DailyClose) error {
	return a.repo.BatchUpsertDailyClose(ctx, objs)
}

func (a *Imp) CreateDailyClose(ctx context.Context, obj *domain.DailyClose) error {
	return a.repo.CreateDailyClose(ctx, obj)
}

func (a *Imp) HasDailyClose(ctx context.Context, date string) (bool, error) {
	return a.repo.HasDailyClose(ctx, date)
}

func (a *Imp) ListDailyClose(
	ctx context.Context,
	arg *domain.ListDailyCloseParams,
) ([]*domain.DailyClose, error) {
	return a.repo.ListDailyClose(ctx, arg)
}

func (a *Imp) ListLatestPrice(
	ctx context.Context,
	stockIDs []string,
) ([]*domain.StockPrice, error) {
	return a.repo.ListLatestPrice(ctx, stockIDs)
}
