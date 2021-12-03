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
package dto

import "samwang0723/jarvis/entity"

type ListDailyCloseSearchParams struct {
	StockIDs *[]string `json:"StockIDs,omitempty"`
	Start    string    `json:"Start"`
	End      *string   `json:"End,omitempty"`
}

type ListDailyCloseRequest struct {
	Offset       int                         `json:"Offset"`
	Limit        int                         `json:"Limit"`
	SearchParams *ListDailyCloseSearchParams `json:"SearchParams"`
}

type ListDailyCloseResponse struct {
	Offset     int                  `json:"Offset"`
	Limit      int                  `json:"Limit"`
	TotalCount int                  `json:"TotalCount"`
	Entries    []*entity.DailyClose `json:"Entries"`
}

type ListStockSearchParams struct {
	StockIDs *[]string `json:"StockIDs,omitempty"`
	Country  string    `json:"Country"`
}

type ListStockRequest struct {
	Offset       int                    `json:"Offset"`
	Limit        int                    `json:"Limit"`
	SearchParams *ListStockSearchParams `json:"SearchParams"`
}

type ListStockResponse struct {
	Offset     int             `json:"Offset"`
	Limit      int             `json:"Limit"`
	TotalCount int             `json:"TotalCount"`
	Entries    []*entity.Stock `json:"Entries"`
}

type DownloadRequest struct {
	RewindLimit int `json:"RewindLimit"`
	RateLimit   int `json:"RateLimit"`
}

type CreateStakeConcentrationRequest struct {
	StockID       string  `json:"StockID"`
	Date          string  `json:"Date"`
	SumBuyShares  uint64  `json:"SumBuyShares"`
	SumSellShares uint64  `json:"SumSellShares"`
	AvgBuyPrice   float32 `json:"AvgBuyPrice"`
	AvgSellPrice  float32 `json:"AvgSellPrice"`
}

type CreateStakeConcentrationResponse struct {
	Entry *entity.StakeConcentration `json:"Entry"`
}

type GetStakeConcentrationRequest struct {
	StockID string `json:"StockID"`
}
