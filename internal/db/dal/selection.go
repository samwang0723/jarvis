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
	"math"
	"time"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	minDailyVolume      = 3000000
	highestRangePercent = 0.04
	maxRewind           = -180
)

func (i *dalImpl) ListSelections(ctx context.Context, offset int32, limit int32,
	date string,
) (objs []*entity.Selection, totalCount int64, err error) {
	err = i.db.Raw(`select count(*) from stake_concentration s
			left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = ?)
			where (
				IF(concentration_1 > 0, 1, 0) +
				IF(concentration_5 > 0, 1, 0) +
				IF(concentration_10 > 0, 1, 0) +
				IF(concentration_20 > 0, 1, 0) +
				IF(concentration_60 > 0, 1, 0)
			) >= 4
			and s.exchange_date = ?
			and d.close / d.high <= 1.0
			and d.close / d.open >= 1.04
			and d.trade_shares >= ?`, date, date, minDailyVolume).Scan(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = i.db.Raw(`select s.stock_id, c.name, c.category, s.exchange_date, d.open, d.close, d.high, d.low, d.price_diff,
			s.concentration_1, s.concentration_5, s.concentration_10, s.concentration_20, s.concentration_60
			, floor(d.trade_shares/1000) as volume, floor(t.foreign_trade_shares/1000) as foreignc,
			floor(t.trust_trade_shares/1000) as trust, floor(t.hedging_trade_shares/1000) as hedging,
			floor(t.dealer_trade_shares/1000) as dealer
			from stake_concentration s
			left join stocks c on c.stock_id = s.stock_id
			left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = ?)
			left join three_primary t on (t.stock_id = s.stock_id and t.exchange_date = ?)
			where (
				IF(s.concentration_1 > 0, 1, 0) +
				IF(s.concentration_5 > 0, 1, 0) +
				IF(s.concentration_10 > 0, 1, 0) +
				IF(s.concentration_20 > 0, 1, 0) +
				IF(s.concentration_60 > 0, 1, 0)
			) >= 4
			and s.exchange_date = ?
			and d.close / d.high <= 1.0
			and d.close / d.open >= 1.04
			and d.trade_shares >= ?
			order by s.concentration_1 desc limit ?, ?`, date, date, date, minDailyVolume,
		offset, limit).Scan(&objs).Error
	if err != nil {
		return nil, 0, err
	}

	// doing analysis
	err = i.expectantHighest(objs)
	if err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}

func (i *dalImpl) expectantHighest(objs []*entity.Selection) error {
	selectionMap := make(map[string]*entity.Selection)
	stockIDs := make([]string, 0, len(objs))
	currentDate := ""
	for _, obj := range objs {
		stockIDs = append(stockIDs, obj.StockID)
		selectionMap[obj.StockID] = obj
		if currentDate == "" {
			currentDate = obj.Date
		}
	}

	t := time.Now().AddDate(0, 0, maxRewind)
	startDate := t.Format("20060102")

	type max180days struct {
		StockID    string  `gorm:"column:stock_id"`
		Max180days float32 `gorm:"column:max_180"`
	}

	var max180daysList []*max180days

	err := i.db.Raw(`select stock_id, MAX(close) as max_180
			from daily_closes 
			WHERE stock_id IN (?)
			and exchange_date >= ? and exchange_date < ?
			GROUP BY stock_id`, stockIDs, startDate, currentDate).Scan(&max180daysList).Error
	if err != nil {
		return err
	}

	for _, max := range max180daysList {
		ref := selectionMap[max.StockID]
		ref.ExpectantHighest = math.Abs(1.0-float64(ref.Close/max.Max180days)) <= highestRangePercent
	}

	return nil
}
