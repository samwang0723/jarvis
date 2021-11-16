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
	"fmt"
	"reflect"
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/services/convert"
)

func (s *serviceImpl) BatchCreateDailyClose(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.DailyClose
	dailyCloses := []*entity.DailyClose{}
	for _, v := range *objs {
		if val, ok := v.(*entity.DailyClose); ok {
			dailyCloses = append(dailyCloses, val)
		} else {
			return fmt.Errorf("cannot cast interface to *dto.DailyClose: %v\n", reflect.TypeOf(v).Elem())
		}
	}

	// prevent Prepared statement contains too many placeholders
	// maxmium 1,000 records
	var err error
	l := len(dailyCloses)
	for i := 0; i < l; i += idal.MaxRow {
		end := i + idal.MaxRow
		if end > l {
			end = l
		}
		err = s.dal.BatchCreateDailyClose(ctx, dailyCloses[i:end])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serviceImpl) ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error) {
	objs, totalCount, err := s.dal.ListDailyClose(ctx, req.Offset, req.Limit, convert.ListDailyCloseSearchParamsDTOToDAL(req.SearchParams))
	if err != nil {
		return nil, 0, err
	}
	return objs, totalCount, nil
}
