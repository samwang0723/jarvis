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
package dal

import (
	"context"
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/entity"
	"strings"

	"gorm.io/gorm"
)

func (i *dalImpl) CreateStock(ctx context.Context, obj *entity.Stock) error {
	err := i.db.Create(obj).Error
	return err
}

func (i *dalImpl) UpdateStock(ctx context.Context, obj *entity.Stock) error {
	err := i.db.Unscoped().Model(&entity.Stock{}).Save(obj).Error
	return err
}

func (i *dalImpl) DeleteStockByID(ctx context.Context, id entity.ID) error {
	err := i.db.Delete(&entity.Stock{}, id).Error
	return err
}

func (i *dalImpl) GetStockByStockID(ctx context.Context, stockID string) (*entity.Stock, error) {
	res := &entity.Stock{}
	if err := i.db.First(res, "stock_id = ?", stockID).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (i *dalImpl) ListStock(ctx context.Context, offset int, limit int,
	searchParams *idal.ListStockSearchParams) (objs []*entity.Stock, totalCount int64, err error) {
	query := i.db.Model(&entity.Stock{})
	query = buildQueryFromListStockSearchParams(query, searchParams)
	err = query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&objs).Error
	if err != nil {
		return nil, 0, err
	}
	return objs, totalCount, nil
}

func buildQueryFromListStockSearchParams(query *gorm.DB, params *idal.ListStockSearchParams) *gorm.DB {
	if params == nil {
		return query
	}
	if len(params.Country) > 0 {
		query = query.Where("country = ?", params.Country)
	}
	if params.StockIDs != nil {
		query = query.Where("stock_id IN (" + strings.Join(*params.StockIDs, ",") + ")")
	}
	return query
}
