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

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/services/convert"
)

func (s *serviceImpl) BatchUpsertThreePrimary(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.ThreePrimary
	threePrimary := []*entity.ThreePrimary{}
	for _, v := range *objs {
		if val, ok := v.(*entity.ThreePrimary); ok {
			threePrimary = append(threePrimary, val)
		} else {
			return fmt.Errorf("cannot cast interface to *dto.ThreePrimary: %v", reflect.TypeOf(v).Elem())
		}
	}

	return s.dal.BatchUpsertThreePrimary(ctx, threePrimary)
}

func (s *serviceImpl) ListThreePrimary(
	ctx context.Context,
	req *dto.ListThreePrimaryRequest,
) ([]*entity.ThreePrimary, int64, error) {
	objs, totalCount, err := s.dal.ListThreePrimary(
		ctx,
		req.Offset,
		req.Limit,
		convert.ListThreePrimarySearchParamsDTOToDAL(req.SearchParams),
	)
	if err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}
