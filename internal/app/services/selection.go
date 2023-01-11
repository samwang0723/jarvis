// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/cache"
)

const (
	proxyURI                   = "https://api.webscrapingapi.com/v1?api_key=mIuUQw7mBA9hngYdkxYOkKrLtvVjH7Hd&url=%s"
	realTimePriceURI           = "https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=%s"
	realTimeMonitoringKey      = "real_time_monitoring_keys"
	defaultCacheExpire         = 7 * 24 * time.Hour
	defaultRealtimeCacheExpire = 5 * time.Minute
	defaultHTTPTimeout         = 10 * time.Second
	fixedTimeLength            = 2
	rateLimit                  = 2 * time.Second
)

//nolint:nolintlint, gochecknoglobals, gosec
var (
	defaultHTTPClient = &http.Client{
		Timeout: defaultHTTPTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func (s *serviceImpl) ListSelections(ctx context.Context,
	req *dto.ListSelectionRequest,
) (objs []*entity.Selection, err error) {
	// if date is today and time <= 21:00, doing realtime parsing
	// otherwise if date is today and time >= 21:00, doing database calculation
	today := getToday()
	if req.Date != today[0] || (req.Date == today[0] && today[1] >= "2100") {
		objs, err = s.dal.ListSelections(ctx, req.Date, req.Strict)
		if err != nil {
			return nil, err
		}
	}

	return objs, nil
}

func (s *serviceImpl) PresetRealTimeKeys(ctx context.Context) error {
	keys, err := s.dal.GetRealTimeMonitoringKeys(ctx)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s:%s", realTimeMonitoringKey, getToday()[0])
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

func (s *serviceImpl) RetrieveRealTimePrice(ctx context.Context) error {
	keys, err := s.cache.SMembers(ctx, getRedisKey())
	if err != nil {
		return err
	}

	for _, key := range keys {
		go func(ctx context.Context, key string, logger *zerolog.Logger, ca cache.Redis) {
			uri := fmt.Sprintf(realTimePriceURI, key)
			finalURI := fmt.Sprintf(proxyURI, uri)

			req, err := http.NewRequestWithContext(ctx, "GET", finalURI, http.NoBody)
			if err != nil {
				logger.Error().Err(err).Msg("failed to create request")

				return
			}

			req.Header = http.Header{
				"Content-Type": []string{"text/csv;charset=ms950"},
				// It is important to close the connection otherwise fd count will overhead
				"Connection": []string{"close"},
			}
			resp, err := defaultHTTPClient.Do(req)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to do request: %s", key)

				return
			}

			defer resp.Body.Close()

			// Skip payloads for invalid http status codes.
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				logger.Warn().Msgf("response status code is not 2xx: %d, key: %s", resp.StatusCode, key)

				return
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error().Err(err).Msgf("failed to read response body: %s", key)

				return
			}

			rawStr := strings.Trim(string(data), "\n")
			logger.Info().Msg(rawStr)

			// insert temp cache into redis
			redisKey := fmt.Sprintf("%s:%s:temp:%s", realTimeMonitoringKey, getToday()[0], key)
			err = ca.Set(ctx, redisKey, rawStr, defaultRealtimeCacheExpire)
			if err != nil {
				return
			}
		}(ctx, key, s.logger, s.cache)

		time.Sleep(rateLimit)
	}

	return nil
}

func getRedisKey() string {
	return fmt.Sprintf("%s:%s", realTimeMonitoringKey, getToday()[0])
}

func getToday() []string {
	result := make([]string, fixedTimeLength)

	now := time.Now()
	date := now.AddDate(0, 0, 0)
	result[0] = date.Format("20060102")

	hours, minutes, _ := time.Now().Clock()
	result[1] = fmt.Sprintf("%d%02d", hours, minutes)

	return result
}
