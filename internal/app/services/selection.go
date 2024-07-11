package services

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/cache"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	realTimePriceURI           = "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=%s"
	realTimeMonitoringKey      = "real_time_monitoring_keys"
	defaultCacheExpire         = 7 * 24 * time.Hour
	defaultRealtimeCacheExpire = 24 * time.Hour
	rateLimit                  = 2 * time.Second
	skipHeader                 = "skip_dates"
	closeToHighestToday        = 0.985
	realtimeVolume             = 3000
	minDailyVolume             = 3000000
	minWeeklyVolume            = 1000000
	highestRangePercent        = 0.04
	dailyHighestRangePercent   = 0.96
	yesterday                  = 1
	yesterdayAfterClosed       = 2
	priceMA8                   = 8
	priceMA21                  = 21
	priceMA55                  = 55
	volumeMV5                  = 5
	volumeMV13                 = 13
	volumeMV34                 = 34
	threePrimarySumCount       = 10
	rewindWeek                 = -5
)

//nolint:nolintlint,cyclop,nestif
func (s *serviceImpl) ListSelections(ctx context.Context,
	req *dto.ListSelectionRequest,
) (objs []*domain.Selection, err error) {
	today := helper.Today()
	var latestDate string
	hasData, _ := s.dal.HasStakeConcentration(ctx, today)
	if !hasData {
		latestDate, _ = s.dal.GetStakeConcentrationLatestDataPoint(ctx)
	} else {
		latestDate = today
	}

	if req.Date != today || latestDate != "" {
		selections, err := s.dal.ListSelections(ctx, req.Date, req.Strict)
		if err != nil {
			s.logger.Error().Err(err).Msg("list selections data record retrival")

			return nil, err
		}
		// doing analysis
		objs, err = s.advancedFiltering(ctx, selections, req.Strict, req.Date)
		if err != nil {
			s.logger.Error().Err(err).Msg("list selections advanced filtering")

			return nil, err
		}
	} else {
		var realtimeList []domain.Realtime
		var res []*domain.Selection

		chips, err := s.getLatestChip(ctx)
		if err != nil {
			s.logger.Error().Err(err).Msg("get latest chip failed")

			return nil, err
		}

		redisRes, realTimeErr := s.getRealtimeParsedData(ctx, req.Date)
		if realTimeErr != nil {
			s.logger.Error().Err(realTimeErr).Msg("get redis cache failed")

			return nil, realTimeErr
		}

		for _, raw := range redisRes {
			if raw == "" {
				continue
			}

			realtime := &domain.Realtime{}
			e := realtime.UnmarshalJSON([]byte(raw))
			if e != nil || realtime.Close == 0.0 {
				s.logger.Error().Err(e).Msg("unmarshal realtime error")

				continue
			}

			realtimeList = append(realtimeList, *realtime)
		}

		for _, realtime := range realtimeList {
			// override realtime data with history record.
			history := chips[realtime.StockID]
			// if its today, check if reach to highest
			if history == nil || (realtime.Close/realtime.High) <= closeToHighestToday || realtime.Volume < realtimeVolume {
				continue
			}

			obj := &domain.Selection{
				StockID:         realtime.StockID,
				Name:            history.Name,
				ExchangeDate:    realtime.Date,
				Category:        history.Category,
				Open:            realtime.Open,
				High:            realtime.High,
				Low:             realtime.Low,
				Close:           realtime.Close,
				Volume:          int(realtime.Volume),
				PriceDiff:       helper.RoundDecimalTwo(realtime.Close - history.Close),
				Concentration1:  history.Concentration1,
				Concentration5:  history.Concentration5,
				Concentration10: history.Concentration10,
				Concentration20: history.Concentration20,
				Concentration60: history.Concentration60,
				Trust:           history.Trust,
				Dealer:          history.Dealer,
				Foreign:         history.Foreign,
				Hedging:         history.Hedging,
			}

			res = append(res, obj)
		}

		objs, err = s.advancedFiltering(ctx, res, req.Strict, req.Date)
		if err != nil {
			s.logger.Error().Err(err).Msg("advanced filtering failed")

			return nil, err
		}
	}

	return objs, nil
}

func (s *serviceImpl) CronjobPresetRealtimeMonitoringKeys(ctx context.Context) error {
	keys, err := s.dal.GetRealTimeMonitoringKeys(ctx)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s:%s", realTimeMonitoringKey, helper.Today())
	err = s.cache.SAdd(ctx, redisKey, keys)
	if err != nil {
		return err
	}

	err = s.cache.SetExpire(ctx, redisKey, time.Now().Add(defaultCacheExpire))
	if err != nil {
		return err
	}

	return nil
}

//nolint:nolintlint,cyclop
func (s *serviceImpl) CrawlingRealTimePrice(ctx context.Context) error {
	keys, err := s.cache.SMembers(ctx, getRealtimeMonitoringKeys())
	if err != nil {
		return err
	}

	err = s.checkHoliday(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		go func(ctx context.Context, key string, logger *zerolog.Logger, ca cache.Redis) {
			uri := fmt.Sprintf(realTimePriceURI, key)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
			if err != nil {
				logger.Error().Err(err).Msg("failed to create request")

				return
			}

			req.Header = http.Header{
				"Content-Type": []string{"text/csv;charset=ms950"},
				// It is important to close the connection otherwise fd count will overhead
				"Connection": []string{"close"},
			}
			resp, err := s.proxyClient.Do(req)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to do request: %s", key)

				return
			}

			defer resp.Body.Close()

			// Skip payloads for invalid http status codes.
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				logger.Warn().
					Msgf("response status code is not 2xx: %d, key: %s", resp.StatusCode, key)

				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to read response body: %s", key)

				return
			}

			raw := string(data)
			rawStr := strings.ReplaceAll(raw, "\n", "")
			rawStr = strings.ReplaceAll(rawStr, "\\\"", "\"")
			if strings.Contains(rawStr, `"z":"-"`) {
				return
			}

			// insert temp cache into redis
			redisKey := fmt.Sprintf("%s:%s:temp:%s", realTimeMonitoringKey, helper.Today(), key)
			err = ca.Set(ctx, redisKey, rawStr, defaultRealtimeCacheExpire)
			if err != nil {
				return
			}
		}(ctx, key, s.logger, s.cache)

		time.Sleep(rateLimit)
	}

	return nil
}

func (s *serviceImpl) getLatestChip(ctx context.Context) (map[string]*domain.Selection, error) {
	m := make(map[string]*domain.Selection)
	// get latest chip from yesterday
	chip, err := s.dal.GetLatestChip(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range chip {
		m[c.StockID] = c
	}

	return m, nil
}

func (s *serviceImpl) advancedFiltering(
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
		pList, err = s.dal.RetrieveDailyCloseHistory(ctx, stockIDs, opts...)
		wg.Done()
	}()

	go func() {
		tList, err = s.dal.RetrieveThreePrimaryHistory(ctx, stockIDs, opts...)
		wg.Done()
	}()

	go func() {
		if len(objs) > 0 {
			highestPriceMap, err = s.dal.GetHighestPrice(
				ctx,
				stockIDs,
				objs[0].ExchangeDate,
				rewindWeek,
			)
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
) map[string]*domain.Analysis {
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

//nolint:nolintlint,gocognit,cyclop
func filter(
	source map[string]*domain.Selection,
	highestPriceMap map[string]float32,
	analysisMap map[string]*domain.Analysis,
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
