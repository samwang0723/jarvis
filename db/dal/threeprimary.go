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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (i *dalImpl) CreateThreePrimary(ctx context.Context, obj *entity.ThreePrimary) error {
	err := i.db.Create(obj).Error
	return err
}

func (i *dalImpl) BatchUpsertThreePrimary(ctx context.Context, objs []*entity.ThreePrimary) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error
	return err
}

func (i *dalImpl) ListThreePrimary(ctx context.Context, offset int32, limit int32,
	searchParams *idal.ListThreePrimarySearchParams) (objs []*entity.ThreePrimary, totalCount int64, err error) {
	query := i.db.Model(&entity.ThreePrimary{})
	query = buildQueryFromListThreePrimarySearchParams(query, searchParams)
	err = query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(int(offset)).Limit(int(limit)).Find(&objs).Error
	if err != nil {
		return nil, 0, err
	}
	return objs, totalCount, nil
}

func buildQueryFromListThreePrimarySearchParams(query *gorm.DB, params *idal.ListThreePrimarySearchParams) *gorm.DB {
	if params == nil {
		return query
	}
	query = query.Where("exchange_date >= ?", params.Start)
	if params.End != nil {
		dateStr := *params.End
		query = query.Where("exchange_date < ?", dateStr)
	}
	query = query.Where("stock_id = ?", params.StockID)
	return query
}
