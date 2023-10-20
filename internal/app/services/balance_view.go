package services

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) GetBalance(ctx context.Context) (obj *entity.BalanceView, err error) {
	obj, err = s.dal.GetBalanceViewByUserID(ctx, s.currentUserID)
	if err != nil {
		sentry.CaptureException(err)

		return nil, err
	}

	return obj, nil
}
