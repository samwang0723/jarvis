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
package handlers

import (
	"context"

	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/services"
)

type IHandler interface {
	ListDailyClose(
		ctx context.Context,
		req *dto.ListDailyCloseRequest,
	) (*dto.ListDailyCloseResponse, error)
	ListStock(ctx context.Context, req *dto.ListStockRequest) (*dto.ListStockResponse, error)
	ListCategories(ctx context.Context) (*dto.ListCategoriesResponse, error)
	ListSelections(
		ctx context.Context,
		req *dto.ListSelectionRequest,
	) (*dto.ListSelectionResponse, error)
	GetStakeConcentration(
		ctx context.Context,
		req *dto.GetStakeConcentrationRequest,
	) (*domain.StakeConcentration, error)
	ListThreePrimary(
		ctx context.Context,
		req *dto.ListThreePrimaryRequest,
	) (*dto.ListThreePrimaryResponse, error)
	ListeningKafkaInput(ctx context.Context)
	CronjobPresetRealtimeMonitoringKeys(ctx context.Context, schedule string) error
	CrawlingRealTimePrice(ctx context.Context, schedule string) error
	ListPickedStocks(ctx context.Context) (*dto.ListPickedStocksResponse, error)
	InsertPickedStocks(
		ctx context.Context,
		req *dto.InsertPickedStocksRequest,
	) (*dto.InsertPickedStocksResponse, error)
	DeletePickedStocks(
		ctx context.Context,
		req *dto.DeletePickedStocksRequest,
	) (*dto.DeletePickedStocksResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) *dto.LoginResponse
	Logout(ctx context.Context) *dto.LogoutResponse
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	GetBalance(ctx context.Context) (*domain.BalanceView, error)
	CreateTransaction(
		ctx context.Context,
		req *dto.CreateTransactionRequest,
	) (*dto.CreateTransactionResponse, error)
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	ListOrders(ctx context.Context, req *dto.ListOrderRequest) (*dto.ListOrderResponse, error)
}

type handlerImpl struct {
	dataService services.IService
	logger      *zerolog.Logger
	jwtSecret   []byte
}

func New(dataService services.IService, logger *zerolog.Logger) IHandler {
	res := &handlerImpl{
		dataService: dataService,
		logger:      logger,
		jwtSecret:   []byte(config.GetCurrentConfig().JwtSecret),
	}

	return res
}
