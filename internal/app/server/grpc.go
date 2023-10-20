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
package server

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	pb "github.com/samwang0723/jarvis/internal/app/pb"
)

func (s *server) ListDailyClose(ctx context.Context,
	req *pb.ListDailyCloseRequest,
) (*pb.ListDailyCloseResponse, error) {
	res, err := s.Handler().ListDailyClose(ctx, dto.ListDailyCloseRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListDailyCloseResponseToPB(res), nil
}

func (s *server) ListStocks(ctx context.Context, req *pb.ListStockRequest) (*pb.ListStockResponse, error) {
	res, err := s.Handler().ListStock(ctx, dto.ListStockRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListStockResponseToPB(res), nil
}

func (s *server) ListCategories(ctx context.Context,
	req *pb.ListCategoriesRequest,
) (*pb.ListCategoriesResponse, error) {
	res, err := s.Handler().ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ListCategoriesResponseToPB(res), nil
}

func (s *server) GetStakeConcentration(ctx context.Context,
	req *pb.GetStakeConcentrationRequest,
) (*pb.GetStakeConcentrationResponse, error) {
	res, err := s.Handler().GetStakeConcentration(ctx, dto.GetStakeConcentrationRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.GetStakeConcentrationResponseToPB(res), nil
}

func (s *server) ListThreePrimary(ctx context.Context,
	req *pb.ListThreePrimaryRequest,
) (*pb.ListThreePrimaryResponse, error) {
	res, err := s.Handler().ListThreePrimary(ctx, dto.ListThreePrimaryRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListThreePrimaryResponseToPB(res), nil
}

func (s *server) ListSelections(ctx context.Context, req *pb.ListSelectionRequest) (*pb.ListSelectionResponse, error) {
	res, err := s.Handler().ListSelections(ctx, dto.ListSelectionRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListSelectionResponseToPB(res), nil
}

func (s *server) ListPickedStocks(
	ctx context.Context,
	req *pb.ListPickedStocksRequest,
) (*pb.ListPickedStocksResponse, error) {
	res, err := s.Handler().ListPickedStocks(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ListPickedStocksResponseToPB(res), nil
}

func (s *server) InsertPickedStocks(
	ctx context.Context,
	req *pb.InsertPickedStocksRequest,
) (*pb.InsertPickedStocksResponse, error) {
	res, err := s.Handler().InsertPickedStocks(ctx, dto.InsertPickedStocksRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.InsertPickedStocksResponseToPB(res), nil
}

func (s *server) DeletePickedStocks(
	ctx context.Context,
	req *pb.DeletePickedStocksRequest,
) (*pb.DeletePickedStocksResponse, error) {
	res, err := s.Handler().DeletePickedStocks(ctx, dto.DeletePickedStocksRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.DeletePickedStocksResponseToPB(res), nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	res, err := s.Handler().CreateUser(ctx, dto.CreateUserRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.CreateUserResponseToPB(res), nil
}

func (s *server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	res, err := s.Handler().ListUsers(ctx, dto.ListUsersRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListUsersResponseToPB(res), nil
}

func (s *server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	res, err := s.Handler().GetBalance(ctx)
	if err != nil {
		return nil, err
	}

	return dto.GetBalanceResponseToPB(res), nil
}

func (s *server) CreateTransaction(
	ctx context.Context,
	req *pb.CreateTransactionRequest,
) (*pb.CreateTransactionResponse, error) {
	res, err := s.Handler().CreateTransaction(ctx, dto.CreateTransactionRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.CreateTransactionResponseToPB(res), nil
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	res, err := s.Handler().CreateOrder(ctx, dto.CreateOrderRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.CreateOrderResponseToPB(res), nil
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	res, err := s.Handler().ListOrders(ctx, dto.ListOrderRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListOrderResponseToPB(res), nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res := s.Handler().Login(ctx, dto.LoginRequestFromPB(req))

	return dto.LoginResponseToPB(res), nil
}
