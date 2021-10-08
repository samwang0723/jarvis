package idal

import (
	"context"
	"samwang0723/jarvis/entity"
)

type ListStockSearchParams struct {
	StockIDs []*string
	Country  *string
}

type IStockDAL interface {
	CreateStock(ctx context.Context, obj *entity.Stock) error
	UpdateStock(ctx context.Context, obj *entity.Stock) error
	DeleteStockByID(ctx context.Context, id entity.ID) error
	ListStock(ctx context.Context, offset int, limit int,
		searchParams *ListStockSearchParams) (objs []*entity.Stock, totalCount int64, err error)
	GetStockByStockID(ctx context.Context, stockID string) (*entity.Stock, error)
}
