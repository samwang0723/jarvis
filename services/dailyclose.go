package services

import (
	"context"
	"fmt"
	"reflect"
	"samwang0723/jarvis/entity"
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
	err := s.dalService.BatchCreateDailyClose(ctx, dailyCloses)

	return err
}
