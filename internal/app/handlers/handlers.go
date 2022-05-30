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
package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/services"
)

type IHandler interface {
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) (*dto.ListDailyCloseResponse, error)
	ListStock(ctx context.Context, req *dto.ListStockRequest) (*dto.ListStockResponse, error)
	ListCategories(ctx context.Context) (*dto.ListCategoriesResponse, error)
	GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error)
	ListThreePrimary(ctx context.Context, req *dto.ListThreePrimaryRequest) (*dto.ListThreePrimaryResponse, error)
	ListeningKafkaInput(ctx context.Context)
}

type handlerImpl struct {
	dataService services.IService
}

func New(dataService services.IService) IHandler {
	res := &handlerImpl{
		dataService: dataService,
	}
	return res
}
