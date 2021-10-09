package services

import (
	"context"
	"fmt"
	"reflect"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/services/convert"
)

func (s *serviceImpl) BatchCreateDailyClose(ctx context.Context, objs map[string]interface{}) error {
	// Replicate the value from interface to *entity.DailyClose
	dailyCloses := []*entity.DailyClose{}
	for _, v := range objs {
		if val, ok := v.(*entity.DailyClose); ok {
			dailyCloses = append(dailyCloses, val)
		} else {
			return fmt.Errorf("cannot cast interface to *dto.DailyClose: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	err := s.dal.BatchCreateDailyClose(ctx, dailyCloses)

	return err
}

func (s *serviceImpl) ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error) {
	objs, totalCount, err := s.dal.ListDailyClose(ctx, req.Offset, req.Limit, convert.ListDailyCloseSearchParamsDTOToDAL(req.SearchParams))
	if err != nil {
		return nil, 0, err
	}
	return objs, totalCount, nil
}
