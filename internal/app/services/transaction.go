package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) CreateTransaction(
	ctx context.Context,
	orderType string,
	creditAmount, debitAmount float32,
) error {
	transaction, err := entity.NewTransaction(
		s.currentUserID,
		orderType,
		creditAmount,
		debitAmount,
	)
	if err != nil {
		return err
	}

	transactions := []*entity.Transaction{transaction}

	return s.dal.CreateChainTransactions(ctx, transactions)
}
