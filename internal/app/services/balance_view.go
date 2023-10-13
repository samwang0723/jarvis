package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) GetBalanceViewByUserID(
	ctx context.Context,
	userID uint64,
) (obj *entity.BalanceView, err error) {
	obj, err = s.dal.GetBalanceViewByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
