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
)

const minDailyVolume = 3000000

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

	return objs, totalCount, nil
}
