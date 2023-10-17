package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type ListOrderSearchParams struct {
	StockIDs      *[]string
	ExchangeMonth *string
	Status        *string
	UserID        uint64
}

type IOrderDAL interface {
	ListOpenOrders(ctx context.Context, userID uint64, stockID, orderType string) ([]*entity.Order, error)
	CreateOrder(ctx context.Context, orders []*entity.Order, transactions []*entity.Transaction) error
	ListOrders(
		ctx context.Context,
		offset, limit int32,
		searchParams *ListOrderSearchParams,
	) ([]*entity.Order, int64, error)
}
