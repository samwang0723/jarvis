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

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"

	"gorm.io/gorm/clause"
)

func (i *dalImpl) GetStakeConcentrationByStockID(ctx context.Context, stockID string, date string) (*entity.StakeConcentration, error) {
	res := &entity.StakeConcentration{}
	if err := i.db.First(res, "stock_id = ? and exchange_date = ?", stockID, date).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) BatchUpsertStakeConcentration(ctx context.Context, objs []*entity.StakeConcentration) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error

	return err
}

func (i *dalImpl) GetStakeConcentrationsWithVolumes(ctx context.Context, stockId string, date string) (objs []*entity.CalculationBase, err error) {
	if err = i.db.Raw(`SELECT a.trade_shares, 
				CAST(b.sum_buy_shares AS SIGNED) - CAST(b.sum_sell_shares as SIGNED) as diff, 
				a.exchange_date FROM daily_closes a
			left join stake_concentration b on (a.stock_id, a.exchange_date) = (b.stock_id, b.exchange_date)
			where a.stock_id=? and a.exchange_date <= ? 
			order by a.exchange_date desc limit 60`, stockId, date).Scan(&objs).Error; err != nil {
		return nil, err
	}

	return objs, nil
}
