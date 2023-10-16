package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	taiwanStockQuantity = 1000
	dayTradeTaxRate     = 0.5
	taxRate             = 0.003
	feeRate             = 0.001425
	brokerFeeDiscount   = 0.25
)

func (s *serviceImpl) CreateTransaction(ctx context.Context, obj *entity.Transaction) error {
	return nil
}
