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

type Selection struct {
	StockID          string  `gorm:"column:stock_id" json:"stockId"`
	Name             string  `gorm:"column:name" json:"name"`
	Category         string  `gorm:"column:category" json:"category"`
	Date             string  `gorm:"column:exchange_date" json:"exchangeDate"`
	Concentration_1  float32 `gorm:"column:concentration_1" json:"concentration1"`
	Concentration_5  float32 `gorm:"column:concentration_5" json:"concentration5"`
	Concentration_10 float32 `gorm:"column:concentration_10" json:"concentration10"`
	Concentration_20 float32 `gorm:"column:concentration_20" json:"concentration20"`
	Concentration_60 float32 `gorm:"column:concentration_60" json:"concentration60"`
}
