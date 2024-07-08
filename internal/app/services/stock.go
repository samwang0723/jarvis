package services

import (
	"context"
	"errors"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

var errCannotCastStock = errors.New("cannot cast interface to *dto.Stock")

func (s *serviceImpl) BatchUpsertStocks(ctx context.Context, objs *[]any) error {
	// Replicate the value from interface to *domain.Stock
	stocks := []*domain.Stock{}
	for _, v := range *objs {
		if val, ok := v.(*domain.Stock); ok {
			stocks = append(stocks, val)
		} else {
			return errCannotCastStock
		}
	}

	return s.dal.BatchUpsertStocks(ctx, stocks)
}

func (s *serviceImpl) ListStock(
	ctx context.Context,
	req *dto.ListStockRequest,
) ([]*domain.Stock, int64, error) {
	param := &domain.ListStocksParams{
		Offset:  req.Offset,
		Limit:   req.Limit,
		Country: req.SearchParams.Country,
	}

	if req.SearchParams.Category != nil {
		param.Category = *req.SearchParams.Category
	}

	if req.SearchParams.StockIDs != nil {
		param.StockIDs = *req.SearchParams.StockIDs
		param.FilterByStockID = true
	}

	if req.SearchParams.Name != nil {
		param.Name = *req.SearchParams.Name
	}

	objs, err := s.dal.ListStocks(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return objs, int64(len(objs)), nil
}

func (s *serviceImpl) ListCategories(ctx context.Context) (objs []*string, err error) {
	return s.dal.ListCategories(ctx)
}
