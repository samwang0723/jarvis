package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

func (s *serviceImpl) BatchUpsertDailyClose(ctx context.Context, objs *[]any) error {
	// Replicate the value from interface to *domain.DailyClose
	dailyCloses := []*domain.DailyClose{}
	for _, v := range *objs {
		if val, ok := v.(*domain.DailyClose); ok {
			dailyCloses = append(dailyCloses, val)
		} else {
			return errCannotCastDailyClose
		}
	}

	// print the dailyCloses
	for _, v := range dailyCloses {
		s.logger.Info().Str("component", "service").Msgf("dailyClose: %v", v)
	}
	err := s.dal.BatchUpsertDailyClose(ctx, dailyCloses)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) ListDailyClose(
	ctx context.Context,
	req *dto.ListDailyCloseRequest,
) ([]*domain.DailyClose, int64, error) {
	param := &domain.ListDailyCloseParams{
		Limit:     req.Limit,
		Offset:    req.Offset,
		StartDate: req.SearchParams.Start,
		StockID:   req.SearchParams.StockID,
	}
	if req.SearchParams.End != nil {
		param.EndDate = *req.SearchParams.End
	}

	objs, err := s.dal.ListDailyClose(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	length := len(objs)
	capacity := int(req.Limit)
	if capacity > length {
		capacity = length
	}

	return objs[:capacity], int64(length), nil
}

func (s *serviceImpl) HasDailyClose(ctx context.Context, date string) bool {
	has, _ := s.dal.HasDailyClose(ctx, date)
	return has
}
