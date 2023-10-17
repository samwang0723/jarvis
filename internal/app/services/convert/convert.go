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

package convert

import (
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/database/dal/idal"
)

func ListDailyCloseSearchParamsDTOToDAL(obj *dto.ListDailyCloseSearchParams) *idal.ListDailyCloseSearchParams {
	res := &idal.ListDailyCloseSearchParams{
		Start:   obj.Start,
		StockID: obj.StockID,
	}

	if obj.End != nil {
		res.End = obj.End
	}

	return res
}

func ListThreePrimarySearchParamsDTOToDAL(obj *dto.ListThreePrimarySearchParams) *idal.ListThreePrimarySearchParams {
	res := &idal.ListThreePrimarySearchParams{
		StockID: obj.StockID,
		Start:   obj.Start,
	}

	if obj.End != nil {
		res.End = obj.End
	}

	return res
}

func ListStockSearchParamsDTOToDAL(obj *dto.ListStockSearchParams) *idal.ListStockSearchParams {
	res := &idal.ListStockSearchParams{
		Country: obj.Country,
	}

	if obj.StockIDs != nil {
		res.StockIDs = obj.StockIDs
	}

	if obj.Name != nil {
		res.Name = obj.Name
	}

	if obj.Category != nil {
		res.Category = obj.Category
	}

	return res
}

func ListOrderSearchParamsDTOToDAL(obj *dto.ListOrderSearchParams) *idal.ListOrderSearchParams {
	res := &idal.ListOrderSearchParams{
		UserID: obj.UserID,
	}

	if obj.StockIDs != nil {
		res.StockIDs = obj.StockIDs
	}

	if obj.ExchangeMonth != nil {
		res.ExchangeMonth = obj.ExchangeMonth
	}

	if obj.Status != nil {
		res.Status = obj.Status
	}

	return res
}
