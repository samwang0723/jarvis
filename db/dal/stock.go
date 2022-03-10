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

	"github.com/samwang0723/jarvis/db/dal/idal"
	"github.com/samwang0723/jarvis/entity"

	"gorm.io/gorm/clause"
)

func (i *dalImpl) BatchUpsertStocks(ctx context.Context, objs []*entity.Stock) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error
	return err
}

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

func (i *dalImpl) ListStock(ctx context.Context, offset int32, limit int32,
	searchParams *idal.ListStockSearchParams) (objs []*entity.Stock, totalCount int64, err error) {

	sql := fmt.Sprintf("select count(*) from stocks where %s", buildQueryFromListStockSearchParams(searchParams))
	if err = i.db.Raw(sql).Scan(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	sql = fmt.Sprintf(`select t.* from 
		(select id from stocks where %s order by stock_id limit %d, %d) q 
		join stocks t on t.id = q.id`, buildQueryFromListStockSearchParams(searchParams), offset, limit)
	if err = i.db.Raw(sql).Scan(&objs).Error; err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}

func (i *dalImpl) ListCategories(ctx context.Context) (objs []string, err error) {
	if err = i.db.Raw("SELECT category FROM stocks group by category order by stock_id").Scan(&objs).Error; err != nil {
		return []string{}, err
	}
	return objs, nil
}

func buildQueryFromListStockSearchParams(params *idal.ListStockSearchParams) string {
	query := ""
	if params == nil {
		return query
	}
	if len(params.Country) > 0 {
		query = fmt.Sprintf("country = '%s'", params.Country)
	}
	if params.StockIDs != nil {
		idList := ""
		stockIDs := *params.StockIDs
		for i := 0; i < len(stockIDs); i++ {
			if i > 0 {
				idList += ","
			}
			idList += "'" + stockIDs[i] + "'"
		}
		query = fmt.Sprintf("%s and stock_id IN (%s)", query, idList)
	}
	if params.Name != nil {
		query = query + " and name like '%" + *params.Name + "%'"
	}
	if params.Category != nil {
		query = fmt.Sprintf("%s and category = '%s'", query, *params.Category)
	}

	return query
}
