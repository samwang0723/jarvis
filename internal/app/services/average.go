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

	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/helper"
)

func (s *serviceImpl) GetAverages(ctx context.Context, stockID string, startDate string) (*businessmodel.Average, error) {
	objs, err := s.dal.GetHistoricalDailyCloses(ctx, stockID, startDate)
	if err != nil {
		return nil, err
	}

	res := &businessmodel.Average{
		StockID: stockID,
		MA:      make(map[int]float32),
		MV:      make(map[int]uint64),
	}

	cursor := 0
	volumeSum := uint64(0)
	priceSum := float32(0.0)

	for _, obj := range objs {
		cursor++
		priceSum += obj.Close
		if cursor == 8 || cursor == 21 || cursor == 55 {
			res.MA[cursor] = helper.RoundUpDecimalTwo(priceSum / float32(cursor))
		}

		volumeSum += obj.Volume
		if cursor == 5 || cursor == 13 || cursor == 34 {
			res.MV[cursor] = volumeSum / uint64(cursor)
		}
	}

	return res, nil
}
