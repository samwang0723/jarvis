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
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/services"
)

type IHandler interface {
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) (*dto.ListDailyCloseResponse, error)
	ListStock(ctx context.Context, req *dto.ListStockRequest) (*dto.ListStockResponse, error)
	ListCategories(ctx context.Context) (*dto.ListCategoriesResponse, error)
	ListSelections(ctx context.Context, req *dto.ListSelectionRequest) (*dto.ListSelectionResponse, error)
	GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error)
	ListThreePrimary(ctx context.Context, req *dto.ListThreePrimaryRequest) (*dto.ListThreePrimaryResponse, error)
	ListeningKafkaInput(ctx context.Context)
	CronjobPresetRealtimeMonitoringKeys(ctx context.Context, schedule string) error
	RetrieveRealTimePrice(ctx context.Context, schedule string) error
	ListPickedStocks(ctx context.Context) (*dto.ListPickedStocksResponse, error)
	InsertPickedStocks(ctx context.Context, req *dto.InsertPickedStocksRequest) (*dto.InsertPickedStocksResponse, error)
	DeletePickedStocks(ctx context.Context, req *dto.DeletePickedStocksRequest) (*dto.DeletePickedStocksResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	UpdateBalanceView(ctx context.Context, req *dto.UpdateBalanceViewRequest) (*dto.UpdateBalanceViewResponse, error)
	GetBalanceViewByUserID(ctx context.Context, userID uint64) (*entity.BalanceView, error)
	CreateTransactions(ctx context.Context, req *dto.CreateTransactionsRequest) (*dto.CreateTransactionsResponse, error)
}

type handlerImpl struct {
	dataService services.IService
	logger      *zerolog.Logger
}

func New(dataService services.IService, logger *zerolog.Logger) IHandler {
	res := &handlerImpl{
		dataService: dataService,
		logger:      logger,
	}

	return res
}
