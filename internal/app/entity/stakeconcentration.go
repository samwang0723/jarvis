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
	Model

	StockID          string  `gorm:"column:stock_id"`
	Date             string  `gorm:"column:exchange_date"`
	HiddenField      string  `gorm:"-"` // this field is use to identify which period the SumBuyShares/SumSellShares are
	SumBuyShares     uint64  `gorm:"column:sum_buy_shares"`
	SumSellShares    uint64  `gorm:"column:sum_sell_shares"`
	AvgBuyPrice      float32 `gorm:"column:avg_buy_price"`
	AvgSellPrice     float32 `gorm:"column:avg_sell_price"`
	Concentration_1  float32 `gorm:"column:concentration_1"`
	Concentration_5  float32 `gorm:"column:concentration_5"`
	Concentration_10 float32 `gorm:"column:concentration_10"`
	Concentration_20 float32 `gorm:"column:concentration_20"`
	Concentration_60 float32 `gorm:"column:concentration_60"`
}

func (StakeConcentration) TableName() string {
	return "stake_concentration"
}
