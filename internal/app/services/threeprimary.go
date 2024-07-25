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

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

var errCannotCastThreePrimary = errors.New("cannot cast interface to *dto.ThreePrimary")

func (s *serviceImpl) BatchUpsertThreePrimary(ctx context.Context, objs *[]any) error {
	// Replicate the value from interface to *domain.ThreePrimary
	threePrimary := []*domain.ThreePrimary{}
	for _, v := range *objs {
		if val, ok := v.(*domain.ThreePrimary); ok {
			threePrimary = append(threePrimary, val)
		} else {
			return errCannotCastThreePrimary
		}
	}

	return s.dal.BatchUpsertThreePrimary(ctx, threePrimary)
}

func (s *serviceImpl) ListThreePrimary(
	ctx context.Context,
	req *dto.ListThreePrimaryRequest,
) ([]*domain.ThreePrimary, int64, error) {
	param := &domain.ListThreePrimaryParams{
		Offset:    req.Offset,
		Limit:     req.Limit,
		StockID:   req.SearchParams.StockID,
		StartDate: req.SearchParams.Start,
	}

	if req.SearchParams.End != nil {
		param.EndDate = *req.SearchParams.End
	}
	objs, err := s.dal.ListThreePrimary(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return objs, int64(len(objs)), nil
}
