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
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"

	"gorm.io/gorm/clause"
)

const maxAverageLimit = 55 * 3

func (i *dalImpl) CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error {
	err := i.db.Create(obj).Error

	return err
}

func (i *dalImpl) BatchUpsertDailyClose(ctx context.Context, objs []*entity.DailyClose) error {
	err := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(&objs, idal.MaxRow).Error

	return err
}

func (i *dalImpl) HasDailyClose(ctx context.Context, date string) bool {
	res := []string{}
	if err := i.db.Raw(`select stock_id from daily_closes 
				where exchange_date = ? limit 1`, date).Scan(&res).Error; err != nil {
		return false
	}

	return len(res) > 0
}

func (i *dalImpl) UpdateDailyCloseAnalysis(ctx context.Context, obj *entity.DailyClose) error {
	err := i.db.Exec(`update daily_closes set 
			half_year_high = ?, average_fivedays_volume = ?, above_all_ma = ?`,
		obj.HalfYearHigh, obj.AverageFivedaysVolume, obj.AboveAllMA).Error

	return err
}

func (i *dalImpl) ListDailyClose(ctx context.Context, offset int32, limit int32,
	searchParams *idal.ListDailyCloseSearchParams,
) (objs []*entity.DailyClose, totalCount int64, err error) {
	sql := fmt.Sprintf("select count(*) from daily_closes where %s",
		buildQueryFromListDailyCloseSearchParams(searchParams))

	err = i.db.Raw(sql).Scan(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	sql = fmt.Sprintf(`select t.id, t.stock_id, t.exchange_date, t.transactions, 
			floor(t.trade_shares/1000) as trade_shares, floor(t.turnover/1000) as turnover, 
			t.open, t.high, t.close, t.low, t.price_diff, t.created_at, t.updated_at, t.deleted_at from
			(select id from daily_closes where %s order by exchange_date desc limit %d, %d) q
			join daily_closes t on t.id = q.id`,
		buildQueryFromListDailyCloseSearchParams(searchParams), offset, limit+maxAverageLimit)

	err = i.db.Raw(sql).Scan(&objs).Error
	if err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}

func buildQueryFromListDailyCloseSearchParams(params *idal.ListDailyCloseSearchParams) string {
	if params == nil {
		return ""
	}
	query := fmt.Sprintf("exchange_date >= '%s' and stock_id = '%s'", params.Start, params.StockID)
	if params.End != nil {
		dateStr := *params.End
		query = fmt.Sprintf("%s and exchange_date < '%s'", query, dateStr)
	}

	return query
}
