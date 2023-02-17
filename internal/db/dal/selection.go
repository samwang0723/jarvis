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
	"math"
	"sort"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	minDailyVolume           = 3000000
	minWeeklyVolume          = 1000000
	highestRangePercent      = 0.04
	dailyHighestRangePercent = 0.97
	yesterday                = 1
	yesterdayAfterClosed     = 2
	priceMA8                 = 8
	priceMA21                = 21
	priceMA55                = 55
	volumeMV5                = 5
	volumeMV13               = 13
	volumeMV34               = 34
	threePrimarySumCount     = 10
	percent                  = -100
)

type price struct {
	StockID string  `gorm:"column:stock_id"`
	Date    string  `gorm:"column:exchange_date"`
	Close   float32 `gorm:"column:close"`
	Volume  uint64  `gorm:"column:trade_shares"`
}

type threePrimary struct {
	StockID string `gorm:"column:stock_id"`
	Date    string `gorm:"column:exchange_date"`
	Foreign int64  `gorm:"column:foreign_trade_shares"`
	Trust   int64  `gorm:"column:trust_trade_shares"`
	Hedging int64  `gorm:"column:hedging_trade_shares"`
	Dealer  int64  `gorm:"column:dealer_trade_shares"`
}

type analysis struct {
	MA8       float32
	MA21      float32
	MA55      float32
	MV5       uint64
	MV13      uint64
	MV34      uint64
	Foreign   int64
	Trust     int64
	Hedging   int64
	Dealer    int64
	LastClose float32
}

type realTimeList struct {
	StockID string `gorm:"column:stock_id"`
	Market  string `gorm:"column:market"`
}

func (i *dalImpl) DataCompletionDate(ctx context.Context, opts ...string) (date string, err error) {
	if len(opts) > 0 {
		err = i.db.Raw(`select exchange_date from stake_concentration 
			where exchange_date = ? limit 1;`, opts[0]).Scan(&date).Error
	} else {
		err = i.db.Raw(`select exchange_date from stake_concentration 
			order by exchange_date desc limit 1;`).Scan(&date).Error
	}
	if err != nil {
		return "", err
	}

	return date, nil
}

func (i *dalImpl) GetLatestChip(ctx context.Context) ([]*entity.Selection, error) {
	date, err := i.DataCompletionDate(ctx)
	if err != nil {
		return nil, err
	}

	var objs []*entity.Selection
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
			and d.trade_shares >= ?
			order by s.stock_id`, date, date, date, minDailyVolume).Scan(&objs).Error
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (i *dalImpl) GetRealTimeMonitoringKeys(ctx context.Context) ([]string, error) {
	date, err := i.DataCompletionDate(ctx)
	if err != nil {
		return nil, err
	}

	var objs []*realTimeList
	err = i.db.Raw(`select s.stock_id, c.market
			from stake_concentration s
			left join stocks c on c.stock_id = s.stock_id
			left join daily_closes d on (d.stock_id = s.stock_id and d.exchange_date = ?)
			where (
				IF(s.concentration_1 > 0, 1, 0) +
				IF(s.concentration_5 > 0, 1, 0) +
				IF(s.concentration_10 > 0, 1, 0) +
				IF(s.concentration_20 > 0, 1, 0) +
				IF(s.concentration_60 > 0, 1, 0)
			) >= 4
			and s.exchange_date = ?
			and d.trade_shares >= ?`, date, date, minDailyVolume).Scan(&objs).Error
	if err != nil {
		return nil, err
	}

	stockSymbols := make([]string, len(objs))
	for idx, obj := range objs {
		stockSymbols[idx] = fmt.Sprintf("%s_%s.tw", obj.Market, obj.StockID)
	}

	return stockSymbols, nil
}

func (i *dalImpl) ListSelectionsBasedOnPickedStocks(
	ctx context.Context,
	pickedStocks []string,
) (objs []*entity.Selection, err error) {
	var date string
	err = i.db.Raw(`select exchange_date from stake_concentration 
			order by exchange_date desc limit 1;`).Scan(&date).Error
	if err != nil {
		return nil, err
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
                        where s.stock_id in (?)
			and s.exchange_date = ?
			order by s.stock_id`, date, date, pickedStocks, date).Scan(&objs).Error
	if err != nil {
		return nil, err
	}

	output, err := i.concentrationBackfill(objs, pickedStocks)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//nolint:nolintlint,cyclop,gocognit,gocyclo
func (i *dalImpl) concentrationBackfill(objs []*entity.Selection, stockIDs []string) ([]*entity.Selection, error) {
	tList, err := i.retrieveThreePrimaryHistory(stockIDs)
	if err != nil {
		return nil, err
	}

	currentStockID := ""
	currentIdx := 0
	currentTrustSum := int64(0)
	currentForeignSum := int64(0)
	for _, t := range tList {
		if currentStockID != t.StockID {
			currentStockID = t.StockID
			currentIdx = 0
			currentTrustSum = 0
			currentForeignSum = 0
		}

		currentIdx++

		currentTrustSum += t.Trust
		currentForeignSum += t.Foreign

		if currentIdx == threePrimarySumCount {
			for _, obj := range objs {
				if obj.StockID == currentStockID {
					obj.Trust10 = int(currentTrustSum)
					obj.Foreign10 = int(currentForeignSum)
					obj.QuoteChange = helper.RoundDecimalTwo((1 - (obj.Close / (obj.Close - obj.PriceDiff))) * percent)
				}
			}
		}
	}

	return objs, nil
}

func (i *dalImpl) ListSelections(
	ctx context.Context,
	date string,
	strict bool,
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
			and d.trade_shares >= ?
			order by s.stock_id`, date, date, date, minDailyVolume).Scan(&objs).Error
	if err != nil {
		return nil, err
	}

	// doing analysis
	output, err := i.AdvancedFiltering(objs, strict)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//nolint:nolintlint,cyclop,gocognit,gocyclo
func (i *dalImpl) AdvancedFiltering(
	objs []*entity.Selection,
	strict bool,
	opts ...string,
) ([]*entity.Selection, error) {
	selectionMap := make(map[string]*entity.Selection)
	stockIDs := make([]string, len(objs))
	for idx, obj := range objs {
		stockIDs[idx] = obj.StockID
		selectionMap[obj.StockID] = obj
	}

	// TODO: make channel requests
	pList, err := i.retrieveDailyCloseHistory(stockIDs, opts...)
	if err != nil {
		return nil, err
	}

	tList, err := i.retrieveThreePrimaryHistory(stockIDs, opts...)
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

		lastClose := yesterdayAfterClosed
		if len(opts) > 0 {
			lastClose = yesterday
		}

		switch currentIdx {
		case lastClose:
			analysisMap[currentStockID].LastClose = p.Close
		case volumeMV5:
			analysisMap[currentStockID].MV5 = currentVolumeSum / volumeMV5
		case volumeMV13:
			analysisMap[currentStockID].MV13 = currentVolumeSum / volumeMV13
		case volumeMV34:
			analysisMap[currentStockID].MV34 = currentVolumeSum / volumeMV34
		case priceMA8:
			analysisMap[currentStockID].MA8 = currentPriceSum / priceMA8
		case priceMA21:
			analysisMap[currentStockID].MA21 = currentPriceSum / priceMA21
		case priceMA55:
			analysisMap[currentStockID].MA55 = currentPriceSum / priceMA55
		}
	}

	// fulfill concentration data
	currentStockID = ""
	currentIdx = 0
	currentTrustSum := int64(0)
	currentForeignSum := int64(0)
	for _, t := range tList {
		if currentStockID != t.StockID {
			currentStockID = t.StockID
			currentIdx = 0
			currentTrustSum = 0
			currentForeignSum = 0
		}

		currentIdx++

		currentTrustSum += t.Trust
		currentForeignSum += t.Foreign

		if currentIdx == threePrimarySumCount {
			analysisMap[currentStockID].Trust = currentTrustSum
			analysisMap[currentStockID].Foreign = currentForeignSum
		}
	}

	// execute filtering
	output := []*entity.Selection{}
	for k, v := range analysisMap {
		ref := selectionMap[k]
		selected := false

		// if today's realtime value and not within max high range, skip
		if len(opts) > 0 && float64(ref.Close/ref.High) < dailyHighestRangePercent {
			continue
		}

		// checking half-year high is closed enough
		// checking volume is above weekly volume (3000)
		// checking MA8, MA21, MA55 is below today's close
		if math.Abs(1.0-float64(ref.Close/highestPriceMap[ref.StockID])) <= highestRangePercent &&
			v.MV5 >= minWeeklyVolume &&
			ref.Close > v.MA8 &&
			ref.Close > v.MA21 &&
			ref.Close > v.MA55 {
			selected = true
		}

		selectedStrict := false
		if strict &&
			v.MV5 > v.MV13 &&
			v.MV13 > v.MV34 &&
			v.MA8 > v.MA21 &&
			v.MA21 > v.MA55 {
			selectedStrict = true
		}

		if (selected && !strict) || (selected && selectedStrict) {
			ref.Trust10 = int(v.Trust)
			ref.Foreign10 = int(v.Foreign)
			ref.QuoteChange = helper.RoundDecimalTwo((1 - (ref.Close / (ref.Close - ref.PriceDiff))) * percent)
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

	var startDate string
	err := i.db.Raw(`select MIN(a.exchange_date) from (select exchange_date from stake_concentration 
		group by exchange_date order by exchange_date desc limit 120) as a;`).Scan(&startDate).Error
	if err != nil {
		return nil, err
	}

	err = i.db.Raw(`select stock_id, max(high) as high from daily_closes where exchange_date >= ?
			and exchange_date < ? and stock_id IN (?) group by stock_id`, startDate, date, stockIDs).Scan(&highest).Error
	if err != nil {
		return nil, err
	}

	for _, h := range highest {
		highestPriceMap[h.StockID] = h.High
	}

	return highestPriceMap, nil
}

func (i *dalImpl) retrieveDailyCloseHistory(stockIDs []string, opts ...string) ([]*price, error) {
	var pList []*price
	var startDate string
	var err error

	err = i.db.Raw(`select MIN(a.exchange_date) from (select exchange_date from stake_concentration 
		group by exchange_date order by exchange_date desc limit 120) as a;`).Scan(&startDate).Error
	if err != nil {
		return nil, err
	}

	if len(opts) > 0 {
		err = i.db.Raw(`select stock_id, exchange_date, close, trade_shares from daily_closes
			where exchange_date >= ? and exchange_date < ? and stock_id IN (?) order by
			stock_id, exchange_date desc`, startDate, opts[0], stockIDs).Scan(&pList).Error
	} else {
		err = i.db.Raw(`select stock_id, exchange_date, close, trade_shares from daily_closes
			where exchange_date >= ? and stock_id IN (?) order by
			stock_id, exchange_date desc`, startDate, stockIDs).Scan(&pList).Error
	}
	if err != nil {
		return nil, err
	}

	return pList, nil
}

func (i *dalImpl) retrieveThreePrimaryHistory(stockIDs []string, opts ...string) ([]*threePrimary, error) {
	var pList []*threePrimary
	var startDate string
	var err error

	err = i.db.Raw(`select MIN(a.exchange_date) from (select exchange_date from stake_concentration 
		group by exchange_date order by exchange_date desc limit 10) as a;`).Scan(&startDate).Error
	if err != nil {
		return nil, err
	}

	if len(opts) > 0 {
		err = i.db.Raw(`select stock_id, exchange_date, floor(foreign_trade_shares/1000) as foreign_trade_shares, 
			floor(trust_trade_shares/1000) as trust_trade_shares, 
			floor(dealer_trade_shares/1000) as dealer_trade_shares, 
			floor(hedging_trade_shares/1000) as hedging_trade_shares
			from three_primary where exchange_date >= ?
			and exchange_date < ? and stock_id IN (?) 
			order by stock_id, exchange_date desc`, startDate, opts[0], stockIDs).Scan(&pList).Error
	} else {
		err = i.db.Raw(`select stock_id, exchange_date, floor(foreign_trade_shares/1000) as foreign_trade_shares, 
			floor(trust_trade_shares/1000) as trust_trade_shares, 
			floor(dealer_trade_shares/1000) as dealer_trade_shares, 
			floor(hedging_trade_shares/1000) as hedging_trade_shares
			from three_primary where exchange_date >= ? 
			and stock_id IN (?) order by stock_id, exchange_date desc`, startDate, stockIDs).Scan(&pList).Error
	}
	if err != nil {
		return nil, err
	}

	return pList, nil
}
