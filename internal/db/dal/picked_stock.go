// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"

	"gorm.io/gorm/clause"
)

func (i *dalImpl) CreatePickedStock(ctx context.Context, obj *entity.PickedStock) error {
	err := i.db.Create(obj).Error

	return err
}

func (i *dalImpl) BatchUpsertPickedStock(ctx context.Context, objs []*entity.PickedStock) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error

	return err
}

func (i *dalImpl) UpdatePickedStock(ctx context.Context, obj *entity.PickedStock) error {
	err := i.db.Unscoped().Model(&entity.PickedStock{}).Save(obj).Error

	return err
}

func (i *dalImpl) DeletePickedStockByID(ctx context.Context, id entity.ID) error {
	err := i.db.Delete(&entity.PickedStock{}, id).Error

	return err
}

func (i *dalImpl) ListPickedStocks(ctx context.Context) (objs []*entity.Selection, err error) {
	var pickedStocks []*entity.PickedStock
	if err := i.db.Find(&pickedStocks).Error; err != nil {
		return nil, err
	}

	var stockIds []string
	for _, pickedStock := range pickedStocks {
		stockIds = append(stockIds, pickedStock.StockID)
	}

	return i.ListSelectionsBasedOnPickedStocks(ctx, stockIds)
}
