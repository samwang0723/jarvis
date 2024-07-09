package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
)

func (s *serviceImpl) CreateTransaction(
	ctx context.Context,
	orderType string,
	creditAmount, debitAmount float32,
) error {
	transaction, err := domain.NewTransaction(
		s.currentUserID,
		orderType,
		creditAmount,
		debitAmount,
	)
	if err != nil {
		return err
	}

	return s.dal.CreateTransaction(ctx, transaction)
}
