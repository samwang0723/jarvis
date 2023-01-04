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
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"

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
	searchParams *idal.ListThreePrimarySearchParams,
) (objs []*entity.ThreePrimary, totalCount int64, err error) {
	sql := fmt.Sprintf("select count(*) from three_primary where %s", buildQueryFromListThreePrimarySearchParams(searchParams))
	if err = i.db.Raw(sql).Scan(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	sql = fmt.Sprintf(`select t.* from
		(select id from three_primary where %s order by exchange_date desc limit %d, %d) q
		join three_primary t on t.id = q.id`, buildQueryFromListThreePrimarySearchParams(searchParams), offset, limit)

	if err = i.db.Raw(sql).Scan(&objs).Error; err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}

func buildQueryFromListThreePrimarySearchParams(params *idal.ListThreePrimarySearchParams) string {
	if params == nil {
		return ""
	}

	query := fmt.Sprintf("stock_id = '%s'", params.StockID)
	query = fmt.Sprintf("%s and exchange_date >= '%s'", query, params.Start)
	if params.End != nil {
		dateStr := *params.End
		query = fmt.Sprintf("%s and exchange_date < '%s'", query, dateStr)
	}

	return query
}
