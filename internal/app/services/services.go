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
	"crypto/tls"
	"net/http"
	"time"

	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/cache"
	"github.com/samwang0723/jarvis/internal/cronjob"
	"github.com/samwang0723/jarvis/internal/database/dal/idal"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
)

type IService interface {
	BatchUpsertDailyClose(ctx context.Context, objs *[]interface{}) error
	ListDailyClose(
		ctx context.Context,
		req *dto.ListDailyCloseRequest,
	) ([]*entity.DailyClose, int64, error)
	HasDailyClose(ctx context.Context, date string) bool
	BatchUpsertStocks(ctx context.Context, objs *[]interface{}) error
	ListStock(ctx context.Context, req *dto.ListStockRequest) ([]*entity.Stock, int64, error)
	ListCategories(ctx context.Context) (objs []string, err error)
	ListSelections(ctx context.Context, req *dto.ListSelectionRequest) ([]*entity.Selection, error)
	BatchUpsertThreePrimary(ctx context.Context, objs *[]interface{}) error
	ListThreePrimary(
		ctx context.Context,
		req *dto.ListThreePrimaryRequest,
	) ([]*entity.ThreePrimary, int64, error)
	GetStakeConcentration(
		ctx context.Context,
		req *dto.GetStakeConcentrationRequest,
	) (*entity.StakeConcentration, error)
	BatchUpsertStakeConcentration(ctx context.Context, objs *[]interface{}) error
	ListeningKafkaInput(ctx context.Context)
	StopKafka() error
	StopRedis() error
	StartCron()
	StopCron()
	AddJob(ctx context.Context, spec string, job func()) error
	CronjobPresetRealtimeMonitoringKeys(ctx context.Context) error
	CrawlingRealTimePrice(ctx context.Context) error
	BatchUpsertPickedStocks(ctx context.Context, objs []*entity.PickedStock) error
	DeletePickedStockByID(ctx context.Context, stockID string) error
	ListPickedStock(ctx context.Context) ([]*entity.Selection, error)
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
	ListUsers(
		ctx context.Context,
		req *dto.ListUsersRequest,
	) (objs []*entity.User, totalCount int64, err error)
	GetUser(ctx context.Context) (obj *entity.User, err error)
	CreateUser(ctx context.Context, obj *entity.User) (err error)
	UpdateUser(ctx context.Context, obj *entity.User) (err error)
	Login(ctx context.Context, email, password string) (obj *entity.User, err error)
	Logout(ctx context.Context) error
	DeleteUser(ctx context.Context) (err error)
	GetUserByEmail(ctx context.Context, email string) (obj *entity.User, err error)
	GetUserByPhone(ctx context.Context, phone string) (obj *entity.User, err error)
	GetBalance(ctx context.Context) (obj *entity.BalanceView, err error)
	CreateTransaction(
		ctx context.Context,
		orderType string,
		creditAmount, debitAmount float32,
	) error
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) error
	ListOrders(
		ctx context.Context,
		req *dto.ListOrderRequest,
	) (objs []*entity.Order, totalCount int64, err error)
	WithUserID(ctx context.Context) IService
}

const (
	defaultHTTPTimeout = 10 * time.Second
)

type serviceImpl struct {
	dal           idal.IDAL
	consumer      ikafka.IKafka
	cache         cache.Redis
	cronjob       cronjob.Cronjob
	logger        *zerolog.Logger
	currentUserID uint64
	proxyClient   *http.Client
}

//nolint:gosec // skip tls verification
func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}

	if impl.proxyClient == nil {
		impl.proxyClient = &http.Client{
			Timeout: defaultHTTPTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	return impl
}

func (s *serviceImpl) WithUserID(ctx context.Context) IService {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get current user id")

		return s
	}

	boundService := *s
	boundService.currentUserID = userID

	return &boundService
}
