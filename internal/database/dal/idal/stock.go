// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type ListStockSearchParams struct {
	StockIDs *[]string
	Name     *string
	Category *string
	Country  string
}

type IStockDAL interface {
	BatchUpsertStocks(ctx context.Context, objs []*entity.Stock) error
	CreateStock(ctx context.Context, obj *entity.Stock) error
	UpdateStock(ctx context.Context, obj *entity.Stock) error
	DeleteStockByID(ctx context.Context, id entity.ID) error
	ListStock(ctx context.Context, offset, limit int32,
		searchParams *ListStockSearchParams) (objs []*entity.Stock, totalCount int64, err error)
	GetStockByStockID(ctx context.Context, stockID string) (*entity.Stock, error)
	ListCategories(ctx context.Context) (objs []string, err error)
}
