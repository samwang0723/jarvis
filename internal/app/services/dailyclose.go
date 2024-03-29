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
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/services/convert"
)

func (s *serviceImpl) BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.DailyClose
	dailyCloses := []*entity.DailyClose{}
	for _, v := range *objs {
		if val, ok := v.(*entity.DailyClose); ok {
			dailyCloses = append(dailyCloses, val)
		} else {
			sentry.CaptureException(errCannotCastDailyClose)

			return errCannotCastDailyClose
		}
	}

	err := s.dal.BatchUpsertDailyClose(ctx, dailyCloses)
	if err != nil {
		sentry.CaptureException(err)

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
		sentry.CaptureException(err)

		return nil, 0, err
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
