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
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/services/convert"
)

func (s *serviceImpl) CreateStakeConcentration(ctx context.Context, req *dto.CreateStakeConcentrationRequest) error {
	obj, err := convert.StakeConcentrationCreateRequestToEntity(req)
	if err != nil {
		return err
	}
	return s.dal.CreateStakeConcentration(ctx, obj)
}

func (s *serviceImpl) GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error) {
	return s.dal.GetStakeConcentrationByStockID(ctx, req.StockID)
}

func (s *serviceImpl) ListBackfillStakeConcentrationStockIDs(ctx context.Context, date string) ([]string, error) {
	return s.dal.ListBackfillStakeConcentrationStockIDs(ctx, date)
}
