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
	"samwang0723/jarvis/cronjob/icronjob"
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/entity"
)

type IService interface {
	StartCron()
	StopCron()
	AddJob(ctx context.Context, spec string, job func()) error
	BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error)
	BatchUpsertStocks(ctx context.Context, objs *[]interface{}) error
	ListStock(ctx context.Context, req *dto.ListStockRequest) ([]*entity.Stock, int64, error)
	CreateStakeConcentration(ctx context.Context, req *dto.CreateStakeConcentrationRequest) error
	GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error)
}

type serviceImpl struct {
	dal     idal.IDAL
	cronjob icronjob.ICronJob
}

func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}
