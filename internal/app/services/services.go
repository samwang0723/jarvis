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

	"github.com/samwang0723/jarvis/internal/app/businessmodel"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
)

type IService interface {
	BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error)
	HasDailyClose(ctx context.Context, date string) bool
	GetAverages(ctx context.Context, stockID string) (*businessmodel.Average, error)
	BatchUpsertStocks(ctx context.Context, objs *[]interface{}) error
	ListStock(ctx context.Context, req *dto.ListStockRequest) ([]*entity.Stock, int64, error)
	ListCategories(ctx context.Context) (objs []string, err error)
	BatchUpsertThreePrimary(ctx context.Context, objs *[]interface{}) error
	ListThreePrimary(ctx context.Context, req *dto.ListThreePrimaryRequest) ([]*entity.ThreePrimary, int64, error)
	GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error)
	BatchUpsertStakeConcentration(ctx context.Context, objs *[]interface{}) error
	ListeningKafkaInput(ctx context.Context)
	StopKafka() error
}

type serviceImpl struct {
	dal      idal.IDAL
	consumer ikafka.IKafka
}

func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}
