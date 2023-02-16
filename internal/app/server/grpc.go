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
	res, err := s.Handler().ListPickedStocks(ctx, dto.ListPickedStocksRequestFromPB(req))
	if err != nil {
		return nil, err
	}

	return dto.ListPickedStocksResponseToPB(res), nil
}
