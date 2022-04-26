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

type DailyClose struct {
	Model

	StockID      string  `gorm:"column:stock_id"`
	Date         string  `gorm:"column:exchange_date"`
	TradedShares uint64  `gorm:"column:trade_shares"`
	Transactions uint64  `gorm:"column:transactions"`
	Turnover     uint64  `gorm:"column:turnover"`
	Open         float32 `gorm:"column:open"`
	Close        float32 `gorm:"column:close"`
	High         float32 `gorm:"column:high"`
	Low          float32 `gorm:"column:low"`
	PriceDiff    float32 `gorm:"column:price_diff"`
}

func (DailyClose) TableName() string {
	return "daily_closes"
}
