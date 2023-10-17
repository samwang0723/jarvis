package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type IOrderDAL interface {
	CreateOrder(ctx context.Context, orderRequest *entity.Order, transactions []*entity.Transaction) error
}
