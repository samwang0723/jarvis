package adapter

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/adapter/sqlc"
	"github.com/samwang0723/jarvis/internal/app/domain"
)

type Adapter interface {
	ListCategories(ctx context.Context) ([]*string, error)
	ListStocks(ctx context.Context, arg *domain.ListStocksParams) ([]*domain.Stock, error)
	ListThreePrimary(
		ctx context.Context,
		arg *domain.ListThreePrimaryParams,
	) ([]*domain.ThreePrimary, error)
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

func (a *Imp) ListCategories(ctx context.Context) ([]*string, error) {
	return a.repo.ListCategories(ctx)
}

func (a *Imp) ListStocks(
	ctx context.Context,
	arg *domain.ListStocksParams,
) ([]*domain.Stock, error) {
	return a.repo.ListStocks(ctx, arg)
}

func (a *Imp) ListThreePrimary(
	ctx context.Context,
	arg *domain.ListThreePrimaryParams,
) ([]*domain.ThreePrimary, error) {
	return a.repo.ListThreePrimary(ctx, arg)
}
