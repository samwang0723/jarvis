package sqlc

import (
	"context"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/helper"
)

func (repo *Repo) BatchUpsertStakeConcentration(
	ctx context.Context,
	objs []*domain.StakeConcentration,
) error {
	return repo.primary().
		BatchUpsertStakeConcentration(ctx, toSqlcBatchUpsertStakeConcentrationParams(objs))
}

func (repo *Repo) GetStakeConcentrationByStockID(
	ctx context.Context,
	stockID,
	date string,
) (*domain.StakeConcentration, error) {
	result, err := repo.primary().
		GetStakeConcentrationByStockID(ctx, &sqlcdb.GetStakeConcentrationByStockIDParams{
			StockID:      stockID,
			ExchangeDate: date,
		})
	if err != nil {
		return nil, err
	}

	return toDomainStakeConcentration(result), nil
}

func (repo *Repo) GetStakeConcentrationsWithVolumes(
	ctx context.Context,
	stockID, date string,
) ([]*domain.CalculationBase, error) {
	result, err := repo.primary().
		GetStakeConcentrationsWithVolumes(ctx, &sqlcdb.GetStakeConcentrationsWithVolumesParams{
			StockID:      stockID,
			ExchangeDate: date,
		})
	if err != nil {
		return nil, err
	}
	return toDomainCalculationBaseList(result), nil
}

func (repo *Repo) HasStakeConcentration(ctx context.Context, exchangeDate string) (bool, error) {
	return repo.primary().HasStakeConcentration(ctx, exchangeDate)
}

func (repo *Repo) GetStakeConcentrationLatestDataPoint(
	ctx context.Context,
) string {
	exchangeDate, err := repo.primary().GetStakeConcentrationLatestDataPoint(ctx)
	if err != nil {
		return ""
	}
	return exchangeDate
}

func toSqlcBatchUpsertStakeConcentrationParams(
	stakeConcentrations []*domain.StakeConcentration,
) *sqlcdb.BatchUpsertStakeConcentrationParams {
	result := &sqlcdb.BatchUpsertStakeConcentrationParams{
		ID:              make([]uuid.UUID, 0, len(stakeConcentrations)),
		StockID:         make([]string, 0, len(stakeConcentrations)),
		ExchangeDate:    make([]string, 0, len(stakeConcentrations)),
		SumBuyShares:    make([]int64, 0, len(stakeConcentrations)),
		SumSellShares:   make([]int64, 0, len(stakeConcentrations)),
		AvgBuyPrice:     make([]decimal.Big, 0, len(stakeConcentrations)),
		AvgSellPrice:    make([]decimal.Big, 0, len(stakeConcentrations)),
		Concentration1:  make([]decimal.Big, 0, len(stakeConcentrations)),
		Concentration5:  make([]decimal.Big, 0, len(stakeConcentrations)),
		Concentration10: make([]decimal.Big, 0, len(stakeConcentrations)),
		Concentration20: make([]decimal.Big, 0, len(stakeConcentrations)),
		Concentration60: make([]decimal.Big, 0, len(stakeConcentrations)),
	}
	for _, sc := range stakeConcentrations {
		result.ID = append(result.ID, sc.ID.ID)
		result.StockID = append(result.StockID, sc.StockID)
		result.ExchangeDate = append(result.ExchangeDate, sc.Date)
		result.SumBuyShares = append(result.SumBuyShares, int64(sc.SumBuyShares))
		result.SumSellShares = append(result.SumSellShares, int64(sc.SumSellShares))
		result.AvgBuyPrice = append(result.AvgBuyPrice, helper.Float32ToDecimal(sc.AvgBuyPrice))
		result.AvgSellPrice = append(result.AvgSellPrice, helper.Float32ToDecimal(sc.AvgSellPrice))
		result.Concentration1 = append(
			result.Concentration1,
			helper.Float32ToDecimal(sc.Concentration1),
		)
		result.Concentration5 = append(
			result.Concentration5,
			helper.Float32ToDecimal(sc.Concentration5),
		)
		result.Concentration10 = append(
			result.Concentration10,
			helper.Float32ToDecimal(sc.Concentration10),
		)
		result.Concentration20 = append(
			result.Concentration20,
			helper.Float32ToDecimal(sc.Concentration20),
		)
		result.Concentration60 = append(
			result.Concentration60,
			helper.Float32ToDecimal(sc.Concentration60),
		)
	}

	return result
}

func toDomainStakeConcentration(
	res *sqlcdb.StakeConcentration,
) *domain.StakeConcentration {
	time := domain.Time{
		CreatedAt: &res.CreatedAt,
		UpdatedAt: &res.UpdatedAt,
	}
	if res.DeletedAt.Valid {
		time.DeletedAt = &res.DeletedAt.Time
	}
	return &domain.StakeConcentration{
		ID: domain.ID{
			ID: res.ID,
		},
		StockID:         res.StockID,
		Date:            res.ExchangeDate,
		SumBuyShares:    uint64(res.SumBuyShares.Int64),
		SumSellShares:   uint64(res.SumSellShares.Int64),
		AvgBuyPrice:     helper.DecimalToFloat32(res.AvgBuyPrice),
		AvgSellPrice:    helper.DecimalToFloat32(res.AvgSellPrice),
		Concentration1:  helper.DecimalToFloat32(res.Concentration1),
		Concentration5:  helper.DecimalToFloat32(res.Concentration5),
		Concentration10: helper.DecimalToFloat32(res.Concentration10),
		Concentration20: helper.DecimalToFloat32(res.Concentration20),
		Concentration60: helper.DecimalToFloat32(res.Concentration60),
		Time:            time,
	}
}

func toDomainCalculationBaseList(
	res []*sqlcdb.GetStakeConcentrationsWithVolumesRow,
) []*domain.CalculationBase {
	result := make([]*domain.CalculationBase, 0, len(res))
	for _, r := range res {
		result = append(result, toDomainCalculationBase(r))
	}
	return result
}

func toDomainCalculationBase(
	res *sqlcdb.GetStakeConcentrationsWithVolumesRow,
) *domain.CalculationBase {
	return &domain.CalculationBase{
		Date:        res.ExchangeDate,
		TradeShares: uint64(res.TradeShares.Int64),
		Diff:        int(res.Diff),
	}
}
