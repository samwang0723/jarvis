package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
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
		objs, err = s.dal.ListSelections(ctx, req.Date, req.Strict)
		if err != nil {
			s.logger.Error().Err(err).Msg("list selections")

			return nil, err
		}
	} else {
		redisRes, err := s.getRealtimeParsedData(ctx, req.Date)
		if err != nil {
			s.logger.Error().Err(err).Msg("get redis cache failed")

			return nil, err
		}

		var realtimeList []domain.Realtime
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

		chips, err := s.getLatestChip(ctx)
		if err != nil {
			s.logger.Error().Err(err).Msg("get latest chip failed")

			return nil, err
		}

		var res []*domain.Selection
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
				Date:            realtime.Date,
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

		objs, err = s.dal.AdvancedFiltering(ctx, res, req.Strict, req.Date)
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
