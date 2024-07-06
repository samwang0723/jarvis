package sqlc

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) BatchUpsertThreePrimary(
	ctx context.Context,
	objs []*domain.ThreePrimary,
) error {
	return repo.primary().BatchUpsertThreePrimary(ctx, toSqlcBatchUpsertThreePrimaryParams(objs))
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

func (repo *Repo) CreateThreePrimary(
	ctx context.Context,
	arg *domain.ThreePrimary,
) error {
	return repo.primary().CreateThreePrimary(ctx, &sqlcdb.CreateThreePrimaryParams{
		StockID:            arg.StockID,
		ExchangeDate:       arg.ExchangeDate,
		ForeignTradeShares: &arg.ForeignTradeShares,
		TrustTradeShares:   &arg.TrustTradeShares,
		DealerTradeShares:  &arg.DealerTradeShares,
		HedgingTradeShares: &arg.HedgingTradeShares,
	})
}

func toSqlcBatchUpsertThreePrimaryParams(
	threePrimary []*domain.ThreePrimary,
) *sqlcdb.BatchUpsertThreePrimaryParams {
	result := &sqlcdb.BatchUpsertThreePrimaryParams{
		StockID:            make([]string, 0, len(threePrimary)),
		ExchangeDate:       make([]string, 0, len(threePrimary)),
		ForeignTradeShares: make([]int64, 0, len(threePrimary)),
		TrustTradeShares:   make([]int64, 0, len(threePrimary)),
		DealerTradeShares:  make([]int64, 0, len(threePrimary)),
		HedgingTradeShares: make([]int64, 0, len(threePrimary)),
	}
	for _, tp := range threePrimary {
		result.StockID = append(result.StockID, tp.StockID)
		result.ExchangeDate = append(result.ExchangeDate, tp.ExchangeDate)
		result.ForeignTradeShares = append(result.ForeignTradeShares, tp.ForeignTradeShares)
		result.TrustTradeShares = append(result.TrustTradeShares, tp.TrustTradeShares)
		result.DealerTradeShares = append(result.DealerTradeShares, tp.DealerTradeShares)
		result.HedgingTradeShares = append(result.HedgingTradeShares, tp.HedgingTradeShares)
	}

	return result
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
		ForeignTradeShares: *tp.ForeignTradeShares,
		TrustTradeShares:   *tp.TrustTradeShares,
		DealerTradeShares:  *tp.DealerTradeShares,
		HedgingTradeShares: *tp.HedgingTradeShares,
		Time: domain.Time{
			CreatedAt: &tp.CreatedAt.Time,
			UpdatedAt: &tp.UpdatedAt.Time,
			DeletedAt: &tp.DeletedAt.Time,
		},
	}
}