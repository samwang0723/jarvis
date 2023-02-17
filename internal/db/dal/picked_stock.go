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
	"errors"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"

	"gorm.io/gorm/clause"
)

var ErrNoPickedStock = errors.New("no picked stock")

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
	pickedStocks := []string{}
	if err := i.db.Raw(`select stock_id from picked_stocks 
                        where deleted_at is null`).Scan(&pickedStocks).Error; err != nil {
		return nil, err
	}

	if len(pickedStocks) == 0 {
		return nil, ErrNoPickedStock
	}

	return i.ListSelectionsBasedOnPickedStocks(ctx, pickedStocks)
}