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

import "github.com/samwang0723/jarvis/internal/app/entity"

type DownloadType int32

//go:generate stringer -type=DownloadType
const (
	DailyClose DownloadType = iota
	ThreePrimary
	Concentration
	StockList
)

type ListDailyCloseSearchParams struct {
	StockIDs *[]string `json:"stockIDs,omitempty"`
	End      *string   `json:"end,omitempty"`
	Start    string    `json:"start"`
}

type ListDailyCloseRequest struct {
	SearchParams *ListDailyCloseSearchParams `json:"searchParams"`
	Offset       int32                       `json:"offset"`
	Limit        int32                       `json:"limit"`
}

type ListDailyCloseResponse struct {
	Entries    []*entity.DailyClose `json:"entries"`
	Offset     int32                `json:"offset"`
	Limit      int32                `json:"limit"`
	TotalCount int64                `json:"totalCount"`
}

type ListThreePrimarySearchParams struct {
	End     *string `json:"end,omitempty"`
	StockID string  `json:"stockID,omitempty"`
	Start   string  `json:"start"`
}

type ListThreePrimaryRequest struct {
	SearchParams *ListThreePrimarySearchParams `json:"searchParams"`
	Offset       int32                         `json:"offset"`
	Limit        int32                         `json:"limit"`
}

type ListThreePrimaryResponse struct {
	Entries    []*entity.ThreePrimary `json:"entries"`
	Offset     int32                  `json:"offset"`
	Limit      int32                  `json:"limit"`
	TotalCount int64                  `json:"totalCount"`
}

type ListStockSearchParams struct {
	StockIDs *[]string `json:"stockIDs,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Category *string   `json:"category,omitempty"`
	Country  string    `json:"country"`
}

type ListStockRequest struct {
	SearchParams *ListStockSearchParams `json:"searchParams"`
	Offset       int32                  `json:"offset"`
	Limit        int32                  `json:"limit"`
}

type ListStockResponse struct {
	Entries    []*entity.Stock `json:"entries"`
	Offset     int32           `json:"offset"`
	Limit      int32           `json:"limit"`
	TotalCount int64           `json:"totalCount"`
}

type ListCategoriesResponse struct {
	Entries []string `json:"entries"`
}

type DownloadRequest struct {
	UTCTimestamp string         `json:"utcTimestamp"`
	Types        []DownloadType `json:"types"`
	RewindLimit  int32          `json:"rewindLimit"`
	RateLimit    int32          `json:"rateLimit"`
}

type CreateStakeConcentrationRequest struct {
	StockID       string  `json:"stockID"`
	Date          string  `json:"date"`
	SumBuyShares  uint64  `json:"sumBuyShares"`
	SumSellShares uint64  `json:"sumSellShares"`
	AvgBuyPrice   float32 `json:"avgBuyPrice"`
	AvgSellPrice  float32 `json:"avgSellPrice"`
}

type CreateStakeConcentrationResponse struct {
	Entry *entity.StakeConcentration `json:"entry"`
}

type GetStakeConcentrationRequest struct {
	StockID string `json:"stockID"`
	Date    string `json:"date"`
}

type StartCronjobRequest struct {
	Schedule string         `json:"schedule"`
	Types    []DownloadType `json:"types"`
}

type StartCronjobResponse struct {
	Error    string `json:"error"`
	Messages string `json:"messages"`
	Code     int32  `json:"code"`
}

type RefreshStakeConcentrationRequest struct {
	Date    string  `json:"date"`
	StockID string  `json:"stockID"`
	Diff    []int32 `json:"diff"`
}

type RefreshStakeConcentrationResponse struct {
	Error    string `json:"error"`
	Messages string `json:"messages"`
	Code     int32  `json:"code"`
}
