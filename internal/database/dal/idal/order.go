package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type IOrderDAL interface {
	ListOpenOrders(ctx context.Context, userID uint64, stockID, orderType string) ([]*entity.Order, error)
	CreateOrder(ctx context.Context, orders []*entity.Order, transactions []*entity.Transaction) error
}
