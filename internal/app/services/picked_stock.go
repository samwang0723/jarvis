// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	RoundDecimalTwo = 100
)

func (s *serviceImpl) BatchUpsertPickedStocks(ctx context.Context, objs []*entity.PickedStock) error {
	return s.dal.BatchUpsertPickedStock(ctx, objs)
}

func (s *serviceImpl) DeletePickedStockByID(ctx context.Context, stockID string) error {
	return s.dal.DeletePickedStockByID(ctx, stockID)
}

//nolint:nolintlint,cyclop,nestif
func (s *serviceImpl) ListPickedStock(ctx context.Context) ([]*entity.Selection, error) {
	objs, err := s.dal.ListPickedStocks(ctx)
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

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

	// if already had latest stock data from exchange or cannot find redis
	// realtime cache, using the latest database record.
	if latestDate >= today || len(redisRes) == 0 {
		return objs, nil
	}

	realtimeList := make(map[string]*businessmodel.Realtime)
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

	for _, obj := range objs {
		// override realtime data with history record.
		realtime, ok := realtimeList[obj.StockID]
		if !ok {
			continue
		}

		obj.PriceDiff = helper.RoundDecimalTwo(realtime.Close - obj.Close)
		obj.QuoteChange = helper.RoundDecimalTwo(obj.PriceDiff / obj.Close * RoundDecimalTwo)
		obj.Open = realtime.Open
		obj.High = realtime.High
		obj.Low = realtime.Low
		obj.Close = realtime.Close
		obj.Volume = int(realtime.Volume)
	}

	return objs, nil
}
