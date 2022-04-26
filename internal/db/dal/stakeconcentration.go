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

func (i *dalImpl) CreateStakeConcentration(ctx context.Context, obj *entity.StakeConcentration) error {
	err := i.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(obj).Error
	return err
}

func (i *dalImpl) GetStakeConcentrationByStockID(ctx context.Context, stockID string, date string) (*entity.StakeConcentration, error) {
	res := &entity.StakeConcentration{}
	if err := i.db.First(res, "stock_id = ? and exchange_date = ?", stockID, date).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (i *dalImpl) BatchUpdateStakeConcentration(ctx context.Context, objs []*entity.StakeConcentration) error {
	err := i.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "stock_id"}, {Name: "exchange_date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"concentration_1",
			"concentration_5",
			"concentration_10",
			"concentration_20",
			"concentration_60",
		}),
	}).CreateInBatches(&objs, idal.MaxRow).Error
	return err
}

func (i *dalImpl) GetStakeConcentrationsWithVolumes(ctx context.Context, stockId string, date string) (objs []*entity.CalculationBase, err error) {
	if err = i.db.Raw(`SELECT a.trade_shares, CAST(b.sum_buy_shares AS SIGNED) - CAST(b.sum_sell_shares as SIGNED) as diff, a.exchange_date FROM daily_closes a
		left join stake_concentration b on (a.stock_id, a.exchange_date) = (b.stock_id, b.exchange_date)
		where a.stock_id=? and a.exchange_date <= ? order by a.exchange_date desc limit 60`, stockId, date).Scan(&objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (i *dalImpl) ListBackfillStakeConcentrationStockIDs(ctx context.Context, date string) ([]string, error) {
	res := []string{}
	// using reference from daily_closes to keep data alignment
	if err := i.db.Raw(`select a.stock_id from daily_closes as a
		left join stake_concentration as b on (a.stock_id, a.exchange_date) = (b.stock_id, b.exchange_date)
		where b.stock_id is null and a.exchange_date = ?`, date).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (i *dalImpl) HasStakeConcentration(ctx context.Context, date string) bool {
	res := []string{}
	if err := i.db.Raw(`select stock_id from stake_concentration where exchange_date = ? limit 1`, date).Scan(&res).Error; err != nil {
		return false
	}
	return len(res) > 0
}
