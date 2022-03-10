package server

import (
	"context"

	"github.com/samwang0723/jarvis/dto"
	pb "github.com/samwang0723/jarvis/pb"
)

func (s *server) ListDailyClose(ctx context.Context, req *pb.ListDailyCloseRequest) (*pb.ListDailyCloseResponse, error) {
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

func (s *server) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	res, err := s.Handler().ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ListCategoriesResponseToPB(res), nil
}

func (s *server) GetStakeConcentration(ctx context.Context, req *pb.GetStakeConcentrationRequest) (*pb.GetStakeConcentrationResponse, error) {
	res, err := s.Handler().GetStakeConcentration(ctx, dto.GetStakeConcentrationRequestFromPB(req))
	if err != nil {
		return nil, err
	}
	return dto.GetStakeConcentrationResponseToPB(res), nil
}
