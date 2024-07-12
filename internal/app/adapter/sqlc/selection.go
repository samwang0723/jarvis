package sqlc

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/helper"
)

func (repo *Repo) LatestStockStatSnapshot(ctx context.Context) ([]*domain.Selection, error) {
	exchangeDate, err := repo.GetStakeConcentrationLatestDataPoint(ctx)
	if err != nil {
		return nil, err
	}

	res, err := repo.primary().LatestStockStatSnapshot(ctx, exchangeDate)
	if err != nil {
		return nil, err
	}
	return domain.ConvertSelectionList(res), nil
}

func (repo *Repo) GetRealTimeMonitoringKeys(ctx context.Context) ([]string, error) {
	exchangeDate, err := repo.GetStakeConcentrationLatestDataPoint(ctx)
	if err != nil {
		return nil, err
	}

	res, err := repo.primary().GetEligibleStocksFromDate(ctx, exchangeDate)
	if err != nil {
		return nil, err
	}
	objs := make([]*domain.RealtimeList, 0, len(res))
	for _, stock := range res {
		objs = append(objs, &domain.RealtimeList{
			StockID: stock.StockID,
			Market:  *stock.Market,
		})
	}

	resPicked, err := repo.primary().GetEligibleStocksFromPicked(ctx)
	if err != nil {
		return nil, err
	}
	picked := make([]*domain.RealtimeList, 0, len(resPicked))
	for _, stock := range resPicked {
		picked = append(picked, &domain.RealtimeList{
			StockID: stock.StockID,
			Market:  *stock.Market,
		})
	}

	var ordered []*domain.RealtimeList
	resOrdered, err := repo.primary().GetEligibleStocksFromOrder(ctx)
	if err != nil {
		return nil, err
	}
	for _, stock := range resOrdered {
		ordered = append(picked, &domain.RealtimeList{
			StockID: stock.StockID,
			Market:  *stock.Market,
		})
	}

	mergedList := domain.Merge(objs, picked)
	mergedList = domain.Merge(mergedList, ordered)

	stockSymbols := make([]string, len(mergedList))
	for idx, obj := range mergedList {
		stockSymbols[idx] = fmt.Sprintf("%s_%s.tw", obj.Market, obj.StockID)
	}

	return stockSymbols, nil
}

func (repo *Repo) ListSelectionsFromPicked(
	ctx context.Context,
	stockIDs []string,
	exchangeDate string,
) ([]*domain.Selection, error) {
	result, err := repo.primary().
		ListSelectionsFromPicked(ctx, &sqlcdb.ListSelectionsFromPickedParams{
			StockIds:     stockIDs,
			ExchangeDate: exchangeDate,
		})
	if err != nil {
		return nil, err
	}

	return domain.ConvertSelectionList(result), nil
}

func (repo *Repo) ListSelections(
	ctx context.Context,
	date string,
	strict bool,
) ([]*domain.Selection, error) {
	sel, err := repo.primary().ListSelections(ctx, date)
	if err != nil {
		return nil, err
	}
	return domain.ConvertSelectionList(sel), nil
}

func (repo *Repo) GetHighestPrice(
	ctx context.Context,
	stockIDs []string,
	date string,
	rewindWeek int,
) (map[string]float32, error) {
	highestPriceMap := make(map[string]float32, len(stockIDs))
	startDate, err := repo.primary().GetStartDate(ctx, date)
	if err != nil {
		return nil, err
	}

	endDate := helper.RewindDate(date, rewindWeek)
	if endDate == "" {
		endDate = date
	}

	highest, err := repo.primary().GetHighestPrice(ctx, &sqlcdb.GetHighestPriceParams{
		StockIds:  stockIDs,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}

	for _, h := range highest {
		highestPriceMap[h.StockID] = float32(h.High)
	}

	return highestPriceMap, nil
}

func (repo *Repo) RetrieveDailyCloseHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.DailyClose, error) {
	var startDate string
	var err error

	startDate, err = repo.primary().GetStartDate(ctx, opts[0])
	if err != nil {
		return nil, err
	}
	endDate := repo.endDate(ctx, opts[0])
	if endDate == "" {
		res, _ := repo.primary().
			RetrieveDailyCloseHistory(ctx, &sqlcdb.RetrieveDailyCloseHistoryParams{
				StartDate: startDate,
				EndDate:   endDate,
				StockIds:  stockIDs,
			})
		return domain.ConvertDailyCloseList(res), nil
	}
	res, _ := repo.primary().
		RetrieveDailyCloseHistoryWithDate(ctx, &sqlcdb.RetrieveDailyCloseHistoryWithDateParams{
			StartDate: startDate,
			EndDate:   endDate,
			StockIds:  stockIDs,
		})
	return domain.ConvertDailyCloseList(res), nil
}

func (repo *Repo) RetrieveThreePrimaryHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.ThreePrimary, error) {
	startDate, err := repo.primary().GetStartDate(ctx, opts[0])
	if err != nil {
		return nil, err
	}

	endDate := repo.endDate(ctx, opts[0])
	if endDate == "" {
		res, _ := repo.primary().
			RetrieveThreePrimaryHistory(ctx, &sqlcdb.RetrieveThreePrimaryHistoryParams{
				StartDate: startDate,
				EndDate:   endDate,
				StockIds:  stockIDs,
			})
		return domain.ConvertThreePrimaryList(res), nil
	}
	res, _ := repo.primary().
		RetrieveThreePrimaryHistoryWithDate(ctx, &sqlcdb.RetrieveThreePrimaryHistoryWithDateParams{
			StartDate: startDate,
			EndDate:   endDate,
			StockIds:  stockIDs,
		})
	return domain.ConvertThreePrimaryList(res), nil
}

func (repo *Repo) endDate(ctx context.Context, date string) string {
	var endDate string
	if date != "" {
		has, _ := repo.HasStakeConcentration(ctx, date)
		if has {
			endDate = date
		}
	} else {
		date, _ := repo.GetStakeConcentrationLatestDataPoint(ctx)
		endDate = date
	}

	return endDate
}
