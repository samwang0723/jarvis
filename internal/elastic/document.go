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
package elastic

// Documents describes a full-spect stock information including price, volumes, concentration, etc
type Document struct {
	StockID   string
	Name      string
	Country   string
	Category  string
	Date      string
	Open      float32
	High      float32
	Low       float32
	Close     float32
	PriceDiff float32

	Concentration_1  float32
	Concentration_5  float32
	Concentration_10 float32
	Concentration_20 float32
	Concentration_60 float32

	AvgBuyPrice   float32
	AvgSellPrice  float32
	SumBuyShares  uint64
	SumSellShares uint64

	ForeignTradeShares int64
	TrustTradeShares   int64
	DealerTradeShares  int64
	HedgingTradeShares int64
}
