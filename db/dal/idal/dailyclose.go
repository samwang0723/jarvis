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
package idal

import (
	"context"
	"samwang0723/jarvis/entity"
)

type ListDailyCloseSearchParams struct {
	StockIDs *[]string
	Start    string
	End      *string // End = nil means it's querying to up-to-date data
}

type IDailyCloseDAL interface {
	CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error
	BatchUpsertDailyClose(ctx context.Context, objs []*entity.DailyClose) error
	ListDailyClose(ctx context.Context, offset int, limit int,
		searchParams *ListDailyCloseSearchParams) (objs []*entity.DailyClose, totalCount int64, err error)
	HasDailyClose(ctx context.Context, date string) bool
}
