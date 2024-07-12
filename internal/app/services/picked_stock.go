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

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	roundDecimalTwo = 100
	percent         = -100
)

func (s *serviceImpl) BatchUpsertPickedStocks(
	ctx context.Context,
	objs []*domain.PickedStock,
) error {
	for _, obj := range objs {
		obj.UserID = s.currentUserID
	}

	return s.dal.CreatePickedStocks(ctx, objs)
}

func (s *serviceImpl) DeletePickedStockByID(ctx context.Context, stockID string) error {
	return s.dal.DeletePickedStock(ctx, s.currentUserID, stockID)
}

//nolint:nolintlint,cyclop,nestif
func (s *serviceImpl) ListPickedStock(ctx context.Context) ([]*domain.Selection, error) {
	pickedStocks, err := s.dal.ListPickedStocks(ctx, s.currentUserID)
	if err != nil {
		return nil, err
	}
	stockIDs := make([]string, 0, len(pickedStocks))
	for _, pickedStock := range pickedStocks {
		stockIDs = append(stockIDs, pickedStock.StockID)
	}

	selections, err := s.dal.ListSelectionsFromPicked(ctx, stockIDs)
	if err != nil {
		return nil, err
	}
	latestDate := s.dal.GetStakeConcentrationLatestDataPoint(ctx)
	selections, err = s.concentrationBackfill(ctx, selections, stockIDs, latestDate)
	if err != nil {
		return nil, err
	}

	for _, obj := range selections {
		obj.QuoteChange = helper.RoundDecimalTwo(
			(1 - (obj.Close / (obj.Close - obj.PriceDiff))) * percent,
		)
	}

	realtimeList := s.fetchRealtimePrice(ctx)
	for _, obj := range selections {
		// override realtime data with history record.
		realtime, ok := realtimeList[obj.StockID]
		if !ok {
			continue
		}

		obj.PriceDiff = helper.RoundDecimalTwo(realtime.Close - obj.Close)
		obj.QuoteChange = helper.RoundDecimalTwo(obj.PriceDiff / obj.Close * roundDecimalTwo)
		obj.Open = realtime.Open
		obj.High = realtime.High
		obj.Low = realtime.Low
		obj.Close = realtime.Close
		obj.Volume = int(realtime.Volume)
		obj.ExchangeDate = realtime.Date
	}

	return selections, nil
}

func (s *serviceImpl) concentrationBackfill(
	ctx context.Context,
	objs []*domain.Selection,
	stockIDs []string,
	date string,
) ([]*domain.Selection, error) {
	tList, err := s.dal.RetrieveThreePrimaryHistory(ctx, stockIDs, date)
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
