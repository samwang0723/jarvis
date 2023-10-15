package services

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) GetBalanceViewByUserID(
	ctx context.Context,
	userID uint64,
) (obj *entity.BalanceView, err error) {
	obj, err = s.dal.GetBalanceViewByUserID(ctx, userID)
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

	return obj, nil
}
