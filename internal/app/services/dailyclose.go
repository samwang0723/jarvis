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
	"errors"

	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/services/convert"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	priceMA8   = 8
	priceMA21  = 21
	priceMA55  = 55
	volumeMV5  = 5
	volumeMV13 = 13
	volumeMV34 = 34
)

var errCannotCastDailyClose = errors.New("cannot cast interface to *dto.DailyClose")

func (s *serviceImpl) BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.DailyClose
	dailyCloses := []*entity.DailyClose{}
	for _, v := range *objs {
		if val, ok := v.(*entity.DailyClose); ok {
			dailyCloses = append(dailyCloses, val)
		} else {
			return errCannotCastDailyClose
		}
	}

	err := s.dal.BatchUpsertDailyClose(ctx, dailyCloses)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListDailyClose(
	ctx context.Context,
	req *dto.ListDailyCloseRequest,
) ([]*entity.DailyClose, int64, error) {
	objs, totalCount, err := s.dal.ListDailyClose(
		ctx,
		req.Offset,
		req.Limit,
		convert.ListDailyCloseSearchParamsDTOToDAL(req.SearchParams),
	)
	if err != nil {
		return nil, 0, err
	}

	// Calculate the average
	for idx, obj := range objs {
		calculateAverage(obj, objs[idx:])
	}

	capacity := int(req.Limit)
	if capacity > len(objs) {
		capacity = len(objs)
	}

	return objs[:capacity], totalCount, nil
}

func (s *serviceImpl) HasDailyClose(ctx context.Context, date string) bool {
	return s.dal.HasDailyClose(ctx, date)
}

func calculateAverage(target *entity.DailyClose, objs []*entity.DailyClose) {
	cache := &businessmodel.Average{
		MA: make(map[int]float32),
		MV: make(map[int]uint64),
	}
	cursor := 0
	volumeSum := uint64(0)
	priceSum := float32(0)

	for _, obj := range objs {
		cursor++
		priceSum += obj.Close
		if cursor == priceMA8 || cursor == priceMA21 || cursor == priceMA55 {
			cache.MA[cursor] = helper.RoundUpDecimalTwo(priceSum / float32(cursor))
		}

		volumeSum += obj.TradedShares
		if cursor == volumeMV5 || cursor == volumeMV13 || cursor == volumeMV34 {
			cache.MV[cursor] = volumeSum / uint64(cursor)
		}
	}

	target.Average = cache
}
