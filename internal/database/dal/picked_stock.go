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
	"github.com/samwang0723/jarvis/internal/database/dal/idal"
	"github.com/samwang0723/jarvis/internal/helper"
	"gorm.io/gorm/clause"
)

var ErrNoPickedStock = errors.New("no picked stock")

func (i *dalImpl) BatchUpsertPickedStock(ctx context.Context, objs []*entity.PickedStock) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error

	return err
}

func (i *dalImpl) DeletePickedStockByID(ctx context.Context, stockID string) error {
	err := i.db.Exec(`update picked_stocks set deleted_at = NOW() where stock_id = ?`, stockID).Error

	return err
}

func (i *dalImpl) ListPickedStocks(ctx context.Context) (objs []*entity.Selection, err error) {
	pickedStocks := []string{}
	if serr := i.db.Raw(`select stock_id from picked_stocks 
                        where deleted_at is null`).Scan(&pickedStocks).Error; serr != nil {
		return nil, serr
	}

	if len(pickedStocks) == 0 {
		return nil, ErrNoPickedStock
	}

	objs, err = i.ListSelectionsBasedOnPickedStocks(ctx, pickedStocks)
	if err != nil {
		return nil, err
	}

	for _, obj := range objs {
		obj.QuoteChange = helper.RoundDecimalTwo((1 - (obj.Close / (obj.Close - obj.PriceDiff))) * percent)
	}

	return
}
