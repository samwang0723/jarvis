package services

import (
	"context"
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/entity"
)

type IService interface {
	BatchCreateDailyClose(ctx context.Context, objs map[string]interface{}) error
	ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) ([]*entity.DailyClose, int64, error)
}

type serviceImpl struct {
	dal idal.IDAL
}

func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}
