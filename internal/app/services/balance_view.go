package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
)

func (s *serviceImpl) GetBalance(ctx context.Context) (obj *domain.BalanceView, err error) {
	obj, err = s.dal.GetBalanceView(ctx, s.currentUserID)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
