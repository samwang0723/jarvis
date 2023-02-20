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

	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) BatchUpsertPickedStocks(ctx context.Context, objs []*entity.PickedStock) error {
	return s.dal.BatchUpsertPickedStock(ctx, objs)
}

func (s *serviceImpl) DeletePickedStockByID(ctx context.Context, stockID string) error {
	return s.dal.DeletePickedStockByID(ctx, stockID)
}

func (s *serviceImpl) ListPickedStock(ctx context.Context) ([]*entity.Selection, error) {
	objs, err := s.dal.ListPickedStocks(ctx)
	if err != nil {
		return nil, err
	}

	return objs, nil
}
