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
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/helper"
	"golang.org/x/xerrors"
)

// Config encapsulates the settings for configuring the redis service.
type RedisConfig struct {
	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
	// Redis master node DNS hostname
	Master string
	// Redis password
	Password string
	// Redis sentinel addresses
	SentinelAddrs []string
}

func (cfg *RedisConfig) validate() error {
	if cfg.Master == "" {
		return xerrors.Errorf("service.redis.validate: failed, reason: invalid master hostname")
	}

	if len(cfg.SentinelAddrs) == 0 {
		return xerrors.Errorf("service.redis.validate: failed, reason: invalid sentinel addresses")
	}

	return nil
}

func (s *serviceImpl) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	if s.cache == nil {
		return nil
	}

	return s.cache.ObtainLock(ctx, key, expire)
}

func (s *serviceImpl) StopRedis() error {
	//nolint: nolintlint,typecheck
	if s.cache == nil {
		return xerrors.Errorf("service.stopRedis: failed, reason: redis is not running")
	}

	if err := s.cache.Close(); err != nil {
		return xerrors.Errorf("service.stopRedis: failed, reason: cannot stop redis %w", err)
	}

	return nil
}

func (s *serviceImpl) fetchRealtimePrice(ctx context.Context) (map[string]*businessmodel.Realtime, error) {
	today := helper.Today()
	latestDate, err := s.dal.DataCompletionDate(ctx, today)
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

	redisRes, err := s.getRealtimeParsedData(ctx, today)
	if err != nil {
		s.logger.Warn().Err(err).Msg("no redis cache record")
	}

	realtimeList := make(map[string]*businessmodel.Realtime)

	// if already had latest stock data from exchange or cannot find redis
	// realtime cache, using the latest database record.
	if latestDate >= today || len(redisRes) == 0 {
		return realtimeList, nil
	}

	for _, raw := range redisRes {
		if raw == "" {
			continue
		}

		realtime := &businessmodel.Realtime{}
		e := realtime.UnmarshalJSON([]byte(raw))
		if e != nil || realtime.Close == 0.0 {
			sentry.CaptureException(e)

			s.logger.Error().Err(e).Msg("unmarshal realtime error")

			continue
		}

		realtimeList[realtime.StockID] = realtime
	}

	return realtimeList, nil
}

func (s *serviceImpl) getRealtimeParsedData(ctx context.Context, date string) ([]string, error) {
	keys, err := s.cache.SMembers(ctx, getRealtimeMonitoringKeys())
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

	mgetKeys := make([]string, len(keys))
	for idx, key := range keys {
		mgetKeys[idx] = fmt.Sprintf("%s:%s:temp:%s", realTimeMonitoringKey, date, key)
	}

	res, err := s.cache.MGet(ctx, mgetKeys...)
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

	return res, nil
}

func (s *serviceImpl) checkHoliday(ctx context.Context) error {
	skipDates, err := s.cache.SMembers(ctx, skipHeader)
	if err != nil {
		sentry.CaptureException(err)

		return err
	}

	for _, date := range skipDates {
		if date == helper.Today() {
			return xerrors.New("skip holiday")
		}
	}

	return nil
}

func getRealtimeMonitoringKeys() string {
	return fmt.Sprintf("%s:%s", realTimeMonitoringKey, helper.Today())
}
