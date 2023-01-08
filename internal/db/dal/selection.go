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
	"sort"
	"time"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	minDailyVolume      = 3000000
	minWeeklyVolume     = 1000000
	highestRangePercent = 0.04
	maxRewind           = -6
	averageRewind       = -3
	priceMA8            = 8
	priceMA21           = 21
	priceMA55           = 55
	volumeMV5           = 5
)

type price struct {
	StockID string  `gorm:"column:stock_id"`
	Date    string  `gorm:"column:exchange_date"`
	Close   float32 `gorm:"column:close"`
	Volume  uint64  `gorm:"column:trade_shares"`
}

type analysis struct {
	MA8  float32
	MA21 float32
	MA55 float32
	MV5  uint64
}

func (i *dalImpl) ListSelections(
	ctx context.Context,
	date string,
) (objs []*entity.Selection, err error) {
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
			and d.close / d.open >= 1.0
			and d.trade_shares >= ?
			order by s.stock_id`, date, date, date, minDailyVolume).Scan(&objs).Error
	if err != nil {
		return nil, err
	}

	// doing analysis
	output, err := i.advancedFiltering(objs)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//nolint:nolintlint, cyclop
func (i *dalImpl) advancedFiltering(objs []*entity.Selection) ([]*entity.Selection, error) {
	selectionMap := make(map[string]*entity.Selection)
	stockIDs := make([]string, len(objs))
	for idx, obj := range objs {
		stockIDs[idx] = obj.StockID
		selectionMap[obj.StockID] = obj
	}

	pList, err := i.retrieveHistory(stockIDs)
	if err != nil {
		return nil, err
	}

	highestPriceMap, err := i.getHighestPrice(stockIDs, objs[0].Date)
	if err != nil {
		return nil, err
	}

	// giving initial capacity to map to increase performance
	analysisMap := make(map[string]*analysis, len(stockIDs))
	currentIdx := 0
	currentStockID := ""
	currentPriceSum := float32(0)
	currentVolumeSum := uint64(0)
	for _, p := range pList {
		if currentStockID != p.StockID {
			currentStockID = p.StockID
			currentIdx = 0
			currentPriceSum = 0
			currentVolumeSum = 0
			analysisMap[currentStockID] = &analysis{}
		}

		currentIdx++

		currentPriceSum += p.Close
		currentVolumeSum += p.Volume

		switch currentIdx {
		case volumeMV5:
			analysisMap[currentStockID].MV5 = currentVolumeSum / volumeMV5
		case priceMA8:
			analysisMap[currentStockID].MA8 = currentPriceSum / priceMA8
		case priceMA21:
			analysisMap[currentStockID].MA21 = currentPriceSum / priceMA21
		case priceMA55:
			analysisMap[currentStockID].MA55 = currentPriceSum / priceMA55
		}
	}

	output := []*entity.Selection{}
	for k, v := range analysisMap {
		ref := selectionMap[k]

		if math.Abs(1.0-float64(ref.Close/highestPriceMap[ref.StockID])) <= highestRangePercent &&
			v.MV5 >= minWeeklyVolume &&
			ref.Close > v.MA8 &&
			ref.Close > v.MA21 &&
			ref.Close > v.MA55 {
			output = append(output, ref)
		}
	}

	//nolint:nolintlint, gocritic
	sort.Slice(output[:], func(i, j int) bool {
		return output[i].StockID < output[j].StockID
	})

	return output, nil
}

func (i *dalImpl) getHighestPrice(stockIDs []string, date string) (map[string]float32, error) {
	highestPriceMap := make(map[string]float32, len(stockIDs))
	type HighPrice struct {
		StockID string  `gorm:"column:stock_id"`
		High    float32 `gorm:"column:high"`
	}
	highest := []*HighPrice{}

	t := time.Now().AddDate(0, maxRewind, 0)
	startDate := t.Format("20060102")

	err := i.db.Raw(`select stock_id, max(high) as high from daily_closes where exchange_date >= ? 
			and exchange_date < ? and stock_id IN (?) group by stock_id`, startDate, date, stockIDs).Scan(&highest).Error
	if err != nil {
		return nil, err
	}

	for _, h := range highest {
		highestPriceMap[h.StockID] = h.High
	}

	return highestPriceMap, nil
}

func (i *dalImpl) retrieveHistory(stockIDs []string) ([]*price, error) {
	// calculate the max start date
	t := time.Now().AddDate(0, averageRewind, 0)
	startDate := t.Format("20060102")

	var pList []*price
	err := i.db.Raw(`select stock_id, exchange_date, close, trade_shares from daily_closes
			where exchange_date >= ? and stock_id IN (?) order by 
			stock_id, exchange_date desc`, startDate, stockIDs).Scan(&pList).Error
	if err != nil {
		return nil, err
	}

	return pList, nil
}