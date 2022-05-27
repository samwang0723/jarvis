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

type Stock struct {
	Model

	StockID  string `gorm:"column:stock_id" json:"stockId"`
	Name     string `gorm:"column:name" json:"name"`
	Country  string `gorm:"column:country" json:"country"`
	Category string `gorm:"column:category" json:"category"`
}

func (Stock) TableName() string {
	return "stocks"
}
