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

import "github.com/samwang0723/jarvis/entity"

type DownloadType int32

//go:generate stringer -type=DownloadType
const (
	DailyClose DownloadType = iota
	ThreePrimary
	Concentration
)

type ListDailyCloseSearchParams struct {
	StockIDs *[]string `json:"stockIDs,omitempty"`
	Start    string    `json:"start"`
	End      *string   `json:"end,omitempty"`
}

type ListDailyCloseRequest struct {
	Offset       int32                       `json:"offset"`
	Limit        int32                       `json:"limit"`
	SearchParams *ListDailyCloseSearchParams `json:"searchParams"`
}

type ListDailyCloseResponse struct {
	Offset     int32                `json:"offset"`
	Limit      int32                `json:"limit"`
	TotalCount int64                `json:"totalCount"`
	Entries    []*entity.DailyClose `json:"entries"`
}

type ListThreePrimarySearchParams struct {
	StockID string  `json:"stockID,omitempty"`
	Start   string  `json:"start"`
	End     *string `json:"end,omitempty"`
}

type ListThreePrimaryRequest struct {
	Offset       int32                         `json:"offset"`
	Limit        int32                         `json:"limit"`
	SearchParams *ListThreePrimarySearchParams `json:"searchParams"`
}

type ListThreePrimaryResponse struct {
	Offset     int32                  `json:"offset"`
	Limit      int32                  `json:"limit"`
	TotalCount int64                  `json:"totalCount"`
	Entries    []*entity.ThreePrimary `json:"entries"`
}

type ListStockSearchParams struct {
	StockIDs *[]string `json:"stockIDs,omitempty"`
	Country  string    `json:"country"`
	Name     *string   `json:"name,omitempty"`
	Category *string   `json:"category,omitempty"`
}

type ListStockRequest struct {
	Offset       int32                  `json:"offset"`
	Limit        int32                  `json:"limit"`
	SearchParams *ListStockSearchParams `json:"searchParams"`
}

type ListStockResponse struct {
	Offset     int32           `json:"offset"`
	Limit      int32           `json:"limit"`
	TotalCount int64           `json:"totalCount"`
	Entries    []*entity.Stock `json:"entries"`
}

type ListCategoriesResponse struct {
	Entries []string `json:"entries"`
}

type DownloadRequest struct {
	UTCTimestamp string         `json:"utcTimestamp"`
	RewindLimit  int32          `json:"rewindLimit"`
	RateLimit    int32          `json:"rateLimit"`
	Types        []DownloadType `json:"types"`
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
	Code     int32  `json:"code"`
	Error    string `json:"error"`
	Messages string `json:"messages"`
}

type RefreshStakeConcentrationRequest struct {
	Date    string  `json:"date"`
	StockID string  `json:"stockID"`
	Diff    []int32 `json:"diff"`
}

type RefreshStakeConcentrationResponse struct {
	Code     int32  `json:"code"`
	Error    string `json:"error"`
	Messages string `json:"messages"`
}
