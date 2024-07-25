package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	realTimePriceURI           = "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=%s"
	realTimeMonitoringKey      = "real_time_monitoring_keys"
	defaultCacheExpire         = 7 * 24 * time.Hour
	defaultRealtimeCacheExpire = 24 * time.Hour
	rateLimit                  = 2 * time.Second
	skipHeader                 = "skip_dates"
)

// Template for real-time monitoring
type RealTimeMonitoringTemplate struct {
	service *serviceImpl
}

func (t *RealTimeMonitoringTemplate) Execute(ctx context.Context) error {
	keys, err := t.service.cache.SMembers(ctx, getRealtimeMonitoringKeys())
	if err != nil {
		return err
	}

	err = t.service.checkHoliday(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		go t.processKey(ctx, key)
		time.Sleep(rateLimit)
	}

	return nil
}

func (t *RealTimeMonitoringTemplate) processKey(ctx context.Context, key string) {
	uri := fmt.Sprintf(realTimePriceURI, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		t.service.logger.Error().Err(err).Msg("failed to create request")
		return
	}

	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
		"Connection":   []string{"close"},
	}
	resp, err := t.service.proxyClient.Do(req)
	if err != nil {
		t.service.logger.Error().Err(err).Msgf("failed to do request: %s", key)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		t.service.logger.Warn().
			Msgf("response status code is not 2xx: %d, key: %s", resp.StatusCode, key)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.service.logger.Error().Err(err).Msgf("failed to read response body: %s", key)
		return
	}

	raw := string(data)
	rawStr := strings.ReplaceAll(raw, "\n", "")
	rawStr = strings.ReplaceAll(rawStr, "\\\"", "\"")
	if strings.Contains(rawStr, `"z":"-"`) {
		return
	}

	redisKey := fmt.Sprintf("%s:%s:temp:%s", realTimeMonitoringKey, helper.Today(), key)
	err = t.service.cache.Set(ctx, redisKey, rawStr, defaultRealtimeCacheExpire)
	if err != nil {
		return
	}
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

func (s *serviceImpl) CrawlingRealTimePrice(ctx context.Context) error {
	template := &RealTimeMonitoringTemplate{service: s}
	return template.Execute(ctx)
}

func (s *serviceImpl) ListRealTimeSelections(
	ctx context.Context,
	req *dto.ListSelectionRequest,
) ([]*domain.Selection, error) {
	statSnapshot, err := s.latestStockStatSnapshot(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("get latest stock stat snapshot failed")
		return nil, err
	}

	redisRes, err := s.getRealtimeParsedData(ctx, req.Date)
	if err != nil {
		s.logger.Error().Err(err).Msg("get redis cache failed")
		return nil, err
	}

	realtimeList := s.parseRealtimeData(redisRes)
	selections := s.mergeRealtimeToSelection(realtimeList, statSnapshot)
	objs, err := s.executeAnalysisEngine(ctx, selections, req.Strict, req.Date)
	if err != nil {
		s.logger.Error().Err(err).Msg("advanced filtering failed")
		return nil, err
	}

	return objs, nil
}

func (s *serviceImpl) parseRealtimeData(redisRes []string) (realtimeList []domain.Realtime) {
	for _, raw := range redisRes {
		if raw == "" {
			continue
		}

		realtime := &domain.Realtime{}
		if err := realtime.UnmarshalJSON([]byte(raw)); err != nil || realtime.Close == 0.0 {
			s.logger.Error().Err(err).Msg("unmarshal realtime error")
			continue
		}

		realtimeList = append(realtimeList, *realtime)
	}

	return realtimeList
}

func (s *serviceImpl) mergeRealtimeToSelection(
	realtimeList []domain.Realtime,
	chips map[string]*domain.Selection,
) (selections []*domain.Selection) {
	for _, realtime := range realtimeList {
		history := chips[realtime.StockID]
		if history == nil || (realtime.Close/realtime.High) <= closeToHighestToday ||
			realtime.Volume < realtimeVolume {
			continue
		}

		volume, _ := helper.Uint64ToInt(realtime.Volume)
		selection := &domain.Selection{
			StockID:         realtime.StockID,
			Name:            history.Name,
			ExchangeDate:    realtime.Date,
			Category:        history.Category,
			Open:            realtime.Open,
			High:            realtime.High,
			Low:             realtime.Low,
			Close:           realtime.Close,
			Volume:          volume,
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

		selections = append(selections, selection)
	}

	return selections
}
