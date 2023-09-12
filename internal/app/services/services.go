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
	"time"

	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/cache"
	"github.com/samwang0723/jarvis/internal/cronjob"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
)

type IService interface {
	BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error)
	HasDailyClose(ctx context.Context, date string) bool
	BatchUpsertStocks(ctx context.Context, objs *[]interface{}) error
	ListStock(ctx context.Context, req *dto.ListStockRequest) ([]*entity.Stock, int64, error)
	ListCategories(ctx context.Context) (objs []string, err error)
	ListSelections(ctx context.Context, req *dto.ListSelectionRequest) ([]*entity.Selection, error)
	BatchUpsertThreePrimary(ctx context.Context, objs *[]interface{}) error
	ListThreePrimary(ctx context.Context, req *dto.ListThreePrimaryRequest) ([]*entity.ThreePrimary, int64, error)
	GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error)
	BatchUpsertStakeConcentration(ctx context.Context, objs *[]interface{}) error
	ListeningKafkaInput(ctx context.Context)
	StopKafka() error
	StopRedis() error
	StartCron()
	StopCron()
	AddJob(ctx context.Context, spec string, job func()) error
	CronjobPresetRealtimeMonitoringKeys(ctx context.Context) error
	RetrieveRealTimePrice(ctx context.Context) error
	BatchUpsertPickedStocks(ctx context.Context, objs []*entity.PickedStock) error
	DeletePickedStockByID(ctx context.Context, stockID string) error
	ListPickedStock(ctx context.Context) ([]*entity.Selection, error)
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
}

type serviceImpl struct {
	dal      idal.IDAL
	consumer ikafka.IKafka
	cache    cache.Redis
	cronjob  cronjob.Cronjob
	logger   *zerolog.Logger
}

func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}

	return impl
}
