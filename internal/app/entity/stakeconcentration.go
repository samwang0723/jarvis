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

package entity

type CalculationBase struct {
	Date        string `gorm:"column:exchange_date"`
	TradeShares uint64 `gorm:"column:trade_shares"`
	Diff        int    `gorm:"column:diff"`
}

type StakeConcentration struct {
	StockID string `gorm:"column:stock_id"         json:"stockId"`
	Date    string `gorm:"column:exchange_date"    json:"exchangeDate"`
	Base
	Diff            []int32 `gorm:"-"                       json:"diff"`
	SumBuyShares    uint64  `gorm:"column:sum_buy_shares"   json:"sumBuyShares"`
	SumSellShares   uint64  `gorm:"column:sum_sell_shares"  json:"sumSellShares"`
	AvgSellPrice    float32 `gorm:"column:avg_sell_price"   json:"avgSellPrice"`
	Concentration1  float32 `gorm:"column:concentration_1"`
	Concentration5  float32 `gorm:"column:concentration_5"`
	Concentration10 float32 `gorm:"column:concentration_10"`
	Concentration20 float32 `gorm:"column:concentration_20"`
	Concentration60 float32 `gorm:"column:concentration_60"`
	AvgBuyPrice     float32 `gorm:"column:avg_buy_price"    json:"avgBuyPrice"`
}

func (StakeConcentration) TableName() string {
	return "stake_concentration"
}
