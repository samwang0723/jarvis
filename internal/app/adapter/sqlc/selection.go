package sqlc

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	minDailyVolume           = 3000000
	minWeeklyVolume          = 1000000
	highestRangePercent      = 0.04
	dailyHighestRangePercent = 0.96
	yesterday                = 1
	yesterdayAfterClosed     = 2
	priceMA8                 = 8
	priceMA21                = 21
	priceMA55                = 55
	volumeMV5                = 5
	volumeMV13               = 13
	volumeMV34               = 34
	threePrimarySumCount     = 10
	percent                  = -100
	rewindWeek               = -5
)

type analysis struct {
	MA8       float32
	MA21      float32
	MA55      float32
	LastClose float32
	MV5       uint64
	MV13      uint64
	MV34      uint64
	Foreign   int64
	Trust     int64
	Hedging   int64
	Dealer    int64
}

func (repo *Repo) GetLatestChip(ctx context.Context) ([]*domain.Selection, error) {
	exchangeDate, err := repo.GetStakeConcentrationLatestDataPoint(ctx)
	if err != nil {
		return nil, err
	}

	res, err := repo.primary().GetLatestChip(ctx, exchangeDate)
	if err != nil {
		return nil, err
	}
	return toDomainSelectionList(res), nil
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

	mergedList := merge(objs, picked)
	mergedList = merge(mergedList, ordered)

	stockSymbols := make([]string, len(mergedList))
	for idx, obj := range mergedList {
		stockSymbols[idx] = fmt.Sprintf("%s_%s.tw", obj.Market, obj.StockID)
	}

	return stockSymbols, nil
}

func (repo *Repo) ListSelectionsFromPicked(
	ctx context.Context,
	userID uuid.UUID,
) ([]*domain.Selection, error) {
	exchangeDate, err := repo.GetStakeConcentrationLatestDataPoint(ctx)
	if err != nil {
		return nil, err
	}

	pickedStocks, err := repo.ListPickedStocks(ctx, userID)
	if err != nil {
		return nil, err
	}
	stockIDs := make([]string, 0, len(pickedStocks))
	for _, pickedStock := range pickedStocks {
		stockIDs = append(stockIDs, pickedStock.StockID)
	}

	result, err := repo.primary().
		ListSelectionsFromPicked(ctx, &sqlcdb.ListSelectionsFromPickedParams{
			StockIds:     stockIDs,
			ExchangeDate: exchangeDate,
		})
	if err != nil {
		return nil, err
	}

	selections := toDomainSelectionList(result)
	output, err := repo.concentrationBackfill(ctx, selections, stockIDs, exchangeDate)
	if err != nil {
		return nil, err
	}

	return output, nil
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
	selections := toDomainSelectionList(sel)

	// doing analysis
	output, err := repo.AdvancedFiltering(ctx, selections, strict, date)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (repo *Repo) AdvancedFiltering(
	ctx context.Context,
	objs []*domain.Selection,
	strict bool,
	opts ...string,
) ([]*domain.Selection, error) {
	selectionMap := make(map[string]*domain.Selection)
	stockIDs := make([]string, len(objs))
	for idx, obj := range objs {
		stockIDs[idx] = obj.StockID
		selectionMap[obj.StockID] = obj
	}

	var wg sync.WaitGroup
	wg.Add(3)

	var pList []*domain.DailyClose
	var tList []*domain.ThreePrimary
	var highestPriceMap map[string]float32
	var err error

	go func() {
		pList, err = repo.retrieveDailyCloseHistory(ctx, stockIDs, opts...)
		wg.Done()
	}()

	go func() {
		tList, err = repo.retrieveThreePrimaryHistory(ctx, stockIDs, opts...)
		wg.Done()
	}()

	go func() {
		if len(objs) > 0 {
			highestPriceMap, err = repo.getHighestPrice(ctx, stockIDs, objs[0].Date)
		}
		wg.Done()
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	// fulfill analysis materials
	analysisMap := mappingMovingAverageConcentration(pList, tList, len(stockIDs), opts...)

	// filtering based on selection conditions
	output := filter(selectionMap, highestPriceMap, analysisMap, strict, opts...)
	sort.Slice(output, func(i, j int) bool {
		return output[i].StockID < output[j].StockID
	})

	return output, nil
}

//nolint:nolintlint,gocognit,cyclop
func mappingMovingAverageConcentration(
	pList []*domain.DailyClose,
	tList []*domain.ThreePrimary,
	size int,
	opts ...string,
) map[string]*analysis {
	analysisMap := make(map[string]*analysis, size)
	currentIdx := 0
	currentPriceSum := float32(0)
	currentVolumeSum := uint64(0)

	for _, p := range pList {
		if _, ok := analysisMap[p.StockID]; !ok {
			currentIdx = 0
			currentPriceSum = 0
			currentVolumeSum = 0

			analysisMap[p.StockID] = &analysis{}
		}

		currentIdx++
		currentPriceSum += p.Close
		currentVolumeSum += uint64(p.TradedShares)

		lastClose := yesterdayAfterClosed
		if len(opts) > 0 {
			lastClose = yesterday
		}

		switch currentIdx {
		case lastClose:
			analysisMap[p.StockID].LastClose = p.Close
		case volumeMV5:
			analysisMap[p.StockID].MV5 = currentVolumeSum / volumeMV5
		case volumeMV13:
			analysisMap[p.StockID].MV13 = currentVolumeSum / volumeMV13
		case volumeMV34:
			analysisMap[p.StockID].MV34 = currentVolumeSum / volumeMV34
		case priceMA8:
			analysisMap[p.StockID].MA8 = currentPriceSum / priceMA8
		case priceMA21:
			analysisMap[p.StockID].MA21 = currentPriceSum / priceMA21
		case priceMA55:
			analysisMap[p.StockID].MA55 = currentPriceSum / priceMA55
		}
	}

	// fulfill concentration data
	currentStockID := ""
	currentIdx = 0
	currentTrustSum := int64(0)
	currentForeignSum := int64(0)
	for _, t := range tList {
		if currentStockID != t.StockID {
			currentStockID = t.StockID
			currentIdx = 0
			currentTrustSum = 0
			currentForeignSum = 0
		}

		currentIdx++
		currentTrustSum += t.TrustTradeShares
		currentForeignSum += t.ForeignTradeShares

		if currentIdx == threePrimarySumCount {
			analysisMap[currentStockID].Trust = currentTrustSum
			analysisMap[currentStockID].Foreign = currentForeignSum
		}
	}

	return analysisMap
}

//nolint:nolintlint,gocognit,cyclop
func filter(
	source map[string]*domain.Selection,
	highestPriceMap map[string]float32,
	analysisMap map[string]*analysis,
	strict bool,
	opts ...string,
) []*domain.Selection {
	output := []*domain.Selection{}

	for k, v := range analysisMap {
		ref := source[k]
		selected := false

		// if today's realtime value and not within max high range, skip
		if len(opts) > 0 && float64(ref.Close/ref.High) < dailyHighestRangePercent {
			continue
		}

		// checking half-year high is closed enough
		// checking volume is above weekly volume (3000)
		// checking MA8, MA21, MA55 is below today's close
		if math.Abs(1.0-float64(ref.Close/highestPriceMap[ref.StockID])) <= highestRangePercent &&
			v.MV5 >= minWeeklyVolume &&
			ref.Close > v.MA8 &&
			ref.Close > v.MA21 &&
			ref.Close > v.MA55 {
			selected = true
		}

		selectedStrict := false
		if strict &&
			v.MV5 > v.MV13 &&
			v.MV13 > v.MV34 &&
			v.MA8 > v.MA21 &&
			v.MA21 > v.MA55 {
			selectedStrict = true
		}

		if (selected && !strict) || (selected && selectedStrict) {
			ref.Trust10 = int(v.Trust)
			ref.Foreign10 = int(v.Foreign)
			ref.QuoteChange = helper.RoundDecimalTwo(
				(1 - (ref.Close / (ref.Close - ref.PriceDiff))) * percent,
			)
			output = append(output, ref)
		}
	}

	return output
}

func (repo *Repo) getHighestPrice(
	ctx context.Context,
	stockIDs []string,
	date string,
) (map[string]float32, error) {
	highestPriceMap := make(map[string]float32, len(stockIDs))
	startDate, err := repo.primary().GetStartDate(ctx)
	if err != nil {
		return nil, err
	}

	endDate := helper.RewindDate(date, rewindWeek)
	if endDate == "" {
		endDate = date
	}

	highest, err := repo.primary().GetHighestPrice(ctx, &sqlcdb.GetHighestPriceParams{
		StockIds:       stockIDs,
		ExchangeDate:   startDate,
		ExchangeDate_2: endDate,
	})
	if err != nil {
		return nil, err
	}

	for _, h := range highest {
		highestPriceMap[h.StockID] = float32(h.High)
	}

	return highestPriceMap, nil
}

//nolint:nolintlint,dupl
func toDomainSelectionList(sel any) []*domain.Selection {
	var result []*domain.Selection

	switch v := sel.(type) {
	case []*sqlcdb.GetLatestChipRow:
		result = make([]*domain.Selection, 0, len(v))
		for _, s := range v {
			result = append(result, &domain.Selection{
				Name:            *s.Name,
				StockID:         s.StockID,
				Category:        s.Category,
				Date:            s.ExchangeDate,
				Open:            float32(s.Open),
				High:            float32(s.High),
				Low:             float32(s.Low),
				Close:           float32(s.Close),
				PriceDiff:       float32(s.PriceDiff),
				Concentration1:  float32(s.Concentration1),
				Concentration5:  float32(s.Concentration5),
				Concentration10: float32(s.Concentration10),
				Concentration20: float32(s.Concentration20),
				Concentration60: float32(s.Concentration60),
				Volume:          int(s.Volume),
				Trust:           int(s.Trust),
				Foreign:         int(s.Foreignc),
				Hedging:         int(s.Hedging),
				Dealer:          int(s.Dealer),
			})
		}
	case []*sqlcdb.ListSelectionsFromPickedRow:
		result = make([]*domain.Selection, 0, len(v))
		for _, s := range v {
			result = append(result, &domain.Selection{
				Name:            *s.Name,
				StockID:         s.StockID,
				Category:        s.Category,
				Date:            s.ExchangeDate,
				Open:            float32(s.Open),
				High:            float32(s.High),
				Low:             float32(s.Low),
				Close:           float32(s.Close),
				PriceDiff:       float32(s.PriceDiff),
				Concentration1:  float32(s.Concentration1),
				Concentration5:  float32(s.Concentration5),
				Concentration10: float32(s.Concentration10),
				Concentration20: float32(s.Concentration20),
				Concentration60: float32(s.Concentration60),
				Volume:          int(s.Volume),
				Trust:           int(s.Trust),
				Foreign:         int(s.Foreignc),
				Hedging:         int(s.Hedging),
				Dealer:          int(s.Dealer),
			})
		}
	case []*sqlcdb.ListSelectionsRow:
		result = make([]*domain.Selection, 0, len(v))
		for _, s := range v {
			result = append(result, &domain.Selection{
				Name:            *s.Name,
				StockID:         s.StockID,
				Category:        s.Category,
				Date:            s.ExchangeDate,
				Open:            float32(s.Open),
				High:            float32(s.High),
				Low:             float32(s.Low),
				Close:           float32(s.Close),
				PriceDiff:       float32(s.PriceDiff),
				Concentration1:  float32(s.Concentration1),
				Concentration5:  float32(s.Concentration5),
				Concentration10: float32(s.Concentration10),
				Concentration20: float32(s.Concentration20),
				Concentration60: float32(s.Concentration60),
				Volume:          int(s.Volume),
				Trust:           int(s.Trust),
				Foreign:         int(s.Foreignc),
				Hedging:         int(s.Hedging),
				Dealer:          int(s.Dealer),
			})
		}
	default:
		// Handle unexpected types
		panic("unsupported type")
	}

	return result
}

func (repo *Repo) getSearchDate(ctx context.Context, date string) string {
	var searchDate string
	if date != "" {
		has, _ := repo.HasStakeConcentration(ctx, date)
		if has {
			searchDate = date
		}
	} else {
		date, _ := repo.GetStakeConcentrationLatestDataPoint(ctx)
		searchDate = date
	}

	return searchDate
}

func (repo *Repo) retrieveDailyCloseHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.DailyClose, error) {
	var startDate string
	var err error

	startDate, err = repo.primary().GetStartDate(ctx)
	if err != nil {
		return nil, err
	}
	searchDate := repo.getSearchDate(ctx, opts[0])
	if searchDate == "" {
		res, _ := repo.primary().
			RetrieveDailyCloseHistory(ctx, &sqlcdb.RetrieveDailyCloseHistoryParams{
				ExchangeDate:   startDate,
				ExchangeDate_2: searchDate,
				StockIds:       stockIDs,
			})
		return toDomainDailyClose(res), nil
	}
	res, _ := repo.primary().
		RetrieveDailyCloseHistoryWithDate(ctx, &sqlcdb.RetrieveDailyCloseHistoryWithDateParams{
			ExchangeDate:   startDate,
			ExchangeDate_2: searchDate,
			StockIds:       stockIDs,
		})
	return toDomainDailCloseWithDate(res), nil
}

func toDomainDailyClose(objs []*sqlcdb.RetrieveDailyCloseHistoryRow) []*domain.DailyClose {
	result := make([]*domain.DailyClose, 0, len(objs))
	for _, obj := range objs {
		result = append(result, &domain.DailyClose{
			StockID:      obj.StockID,
			Date:         obj.ExchangeDate,
			TradedShares: *obj.TradeShares,
			Close:        float32(obj.Close),
		})
	}
	return result
}

func toDomainDailCloseWithDate(
	objs []*sqlcdb.RetrieveDailyCloseHistoryWithDateRow,
) []*domain.DailyClose {
	result := make([]*domain.DailyClose, 0, len(objs))
	for _, obj := range objs {
		result = append(result, &domain.DailyClose{
			StockID:      obj.StockID,
			Date:         obj.ExchangeDate,
			TradedShares: *obj.TradeShares,
			Close:        float32(obj.Close),
		})
	}
	return result
}

func (repo *Repo) retrieveThreePrimaryHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.ThreePrimary, error) {
	startDate, err := repo.primary().GetStartDate(ctx)
	if err != nil {
		return nil, err
	}

	searchDate := repo.getSearchDate(ctx, opts[0])
	if searchDate == "" {
		res, _ := repo.primary().
			RetrieveThreePrimaryHistory(ctx, &sqlcdb.RetrieveThreePrimaryHistoryParams{
				ExchangeDate:   startDate,
				ExchangeDate_2: searchDate,
				StockIds:       stockIDs,
			})
		return toDomainThreePrimary(res), nil
	}
	res, _ := repo.primary().
		RetrieveThreePrimaryHistoryWithDate(ctx, &sqlcdb.RetrieveThreePrimaryHistoryWithDateParams{
			ExchangeDate:   startDate,
			ExchangeDate_2: searchDate,
			StockIds:       stockIDs,
		})
	return toDomainThreePrimaryWithDate(res), nil
}

func toDomainThreePrimary(objs []*sqlcdb.RetrieveThreePrimaryHistoryRow) []*domain.ThreePrimary {
	result := make([]*domain.ThreePrimary, 0, len(objs))
	for _, obj := range objs {
		result = append(result, &domain.ThreePrimary{
			StockID:            obj.StockID,
			ExchangeDate:       obj.ExchangeDate,
			ForeignTradeShares: int64(obj.ForeignTradeShares),
			TrustTradeShares:   int64(obj.TrustTradeShares),
			DealerTradeShares:  int64(obj.DealerTradeShares),
			HedgingTradeShares: int64(obj.HedgingTradeShares),
		})
	}
	return result
}

func toDomainThreePrimaryWithDate(
	objs []*sqlcdb.RetrieveThreePrimaryHistoryWithDateRow,
) []*domain.ThreePrimary {
	result := make([]*domain.ThreePrimary, 0, len(objs))
	for _, obj := range objs {
		result = append(result, &domain.ThreePrimary{
			StockID:            obj.StockID,
			ExchangeDate:       obj.ExchangeDate,
			ForeignTradeShares: int64(obj.ForeignTradeShares),
			TrustTradeShares:   int64(obj.TrustTradeShares),
			DealerTradeShares:  int64(obj.DealerTradeShares),
			HedgingTradeShares: int64(obj.HedgingTradeShares),
		})
	}
	return result
}

func (repo *Repo) concentrationBackfill(
	ctx context.Context,
	objs []*domain.Selection,
	stockIDs []string,
	date string,
) ([]*domain.Selection, error) {
	tList, err := repo.retrieveThreePrimaryHistory(ctx, stockIDs, date)
	if err != nil {
		return nil, err
	}

	currentStockID := ""
	currentIdx := 0
	currentTrustSum := int64(0)
	currentForeignSum := int64(0)
	for _, t := range tList {
		if currentStockID != t.StockID {
			currentStockID = t.StockID
			currentIdx = 0
			currentTrustSum = 0
			currentForeignSum = 0
		}

		currentIdx++

		currentTrustSum += t.TrustTradeShares
		currentForeignSum += t.ForeignTradeShares

		if currentIdx == threePrimarySumCount {
			for _, obj := range objs {
				if obj.StockID == currentStockID {
					obj.Trust10 = int(currentTrustSum)
					obj.Foreign10 = int(currentForeignSum)
					obj.QuoteChange = helper.RoundDecimalTwo(
						(1 - (obj.Close / (obj.Close - obj.PriceDiff))) * percent,
					)
				}
			}
		}
	}

	return objs, nil
}

func merge(objs, picked []*domain.RealtimeList) []*domain.RealtimeList {
	// Create a map to keep track of seen StockIDs
	seen := make(map[string]bool)

	// Iterate over the objs list and add each object to the merged list if its StockID has not been seen before
	var merged []*domain.RealtimeList
	for _, obj := range objs {
		if _, ok := seen[obj.StockID]; !ok {
			merged = append(merged, obj)
			seen[obj.StockID] = true
		}
	}

	// Iterate over the picked list and add each object to the merged list if its StockID has not been seen before
	for _, obj := range picked {
		if _, ok := seen[obj.StockID]; !ok {
			merged = append(merged, obj)
			seen[obj.StockID] = true
		}
	}

	return merged
}
