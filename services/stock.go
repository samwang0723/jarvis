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

package services

import (
	"context"
	"fmt"
	"reflect"
	"samwang0723/jarvis/entity"
)

func (s *serviceImpl) BatchUpsertStocks(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.DailyClose
	stocks := []*entity.Stock{}
	for _, v := range *objs {
		if val, ok := v.(*entity.Stock); ok {
			stocks = append(stocks, val)
		} else {
			return fmt.Errorf("cannot cast interface to *dto.Stock: %v\n", reflect.TypeOf(v).Elem())
		}
	}

	return s.dal.BatchUpsertStocks(ctx, stocks)
}

func (s *serviceImpl) CreateStock(ctx context.Context, obj *entity.Stock) error {
	return s.dal.CreateStock(ctx, obj)
}
