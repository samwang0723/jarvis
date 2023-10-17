package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) CreateTransaction(ctx context.Context, obj *entity.Transaction) error {
	transactions := []*entity.Transaction{obj}

	return s.dal.CreateChainTransactions(ctx, transactions)
}
