package services

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/samwang0723/jarvis/internal/app/domain"
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
	rewindWeek               = -5
)

func (s *serviceImpl) executeAnalysisEngine(
	ctx context.Context,
	objs []*domain.Selection,
	strict bool,
	opts ...string,
) ([]*domain.Selection, error) {
	selectionMap, stockIDs := mapWithIDs(objs)
	pList, tList, highestPriceMap, err := s.aggregateStockStat(ctx, stockIDs, objs, opts...)
	if err != nil {
		return nil, err
	}

	analysisMap := mapMAToConcentration(pList, tList, len(stockIDs), opts...)
	output := filterByCoreLogic(selectionMap, highestPriceMap, analysisMap, strict, opts...)
	sort.Slice(output, func(i, j int) bool {
		return output[i].StockID < output[j].StockID
	})

	return output, nil
}

func mapWithIDs(objs []*domain.Selection) (map[string]*domain.Selection, []string) {
	defer helper.TrackElapsed(time.Now(), "mapWithIDs")

	selectionMap := make(map[string]*domain.Selection)
	stockIDs := make([]string, len(objs))
	for idx, obj := range objs {
		stockIDs[idx] = obj.StockID
		selectionMap[obj.StockID] = obj
	}
	return selectionMap, stockIDs
}

func (s *serviceImpl) aggregateStockStat(
	ctx context.Context,
	stockIDs []string,
	objs []*domain.Selection,
	opts ...string,
) ([]*domain.DailyClose, []*domain.ThreePrimary, map[string]float32, error) {
	defer helper.TrackElapsed(time.Now(), "aggregateStockStat")

	var wg sync.WaitGroup
	wg.Add(3)

	var pList []*domain.DailyClose
	var tList []*domain.ThreePrimary
	var highestPriceMap map[string]float32

	errChan := make(chan error, 3)

	go func() {
		defer helper.TrackElapsed(time.Now(), " ->> RetrieveDailyCloseHistory")
		defer wg.Done()
		var err error
		pList, err = s.dal.RetrieveDailyCloseHistory(ctx, stockIDs, opts...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer helper.TrackElapsed(time.Now(), " ->> RetrieveThreePrimaryHistory")
		defer wg.Done()
		var err error
		tList, err = s.dal.RetrieveThreePrimaryHistory(ctx, stockIDs, opts...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer helper.TrackElapsed(time.Now(), " ->> GetHighestPrice")
		defer wg.Done()
		if len(objs) > 0 {
			var err error
			highestPriceMap, err = s.dal.GetHighestPrice(
				ctx,
				stockIDs,
				objs[0].ExchangeDate,
				rewindWeek,
			)
			if err != nil {
				errChan <- err
			}
		}
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return pList, tList, highestPriceMap, nil
}

func mapMAToConcentration(
	pList []*domain.DailyClose,
	tList []*domain.ThreePrimary,
	size int,
	opts ...string,
) map[string]*domain.Analysis {
	defer helper.TrackElapsed(time.Now(), "mapMAToConcentration")

	analysisMap := make(map[string]*domain.Analysis, size)
	currentIdx := 0
	currentPriceSum := float32(0)
	currentVolumeSum := uint64(0)

	for _, p := range pList {
		if _, ok := analysisMap[p.StockID]; !ok {
			currentIdx = 0
			currentPriceSum = 0
			currentVolumeSum = 0

			analysisMap[p.StockID] = &domain.Analysis{}
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

func filterByCoreLogic(
	source map[string]*domain.Selection,
	highestPriceMap map[string]float32,
	analysisMap map[string]*domain.Analysis,
	strict bool,
	opts ...string,
) []*domain.Selection {
	defer helper.TrackElapsed(time.Now(), "filterByCoreLogic")

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
