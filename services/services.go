package services

import (
	"context"
	"samwang0723/jarvis/db"
	"samwang0723/jarvis/db/dal"
	"samwang0723/jarvis/db/dal/idal"
)

type IService interface {
	BatchCreateDailyClose(ctx context.Context, objs map[string]interface{}) error
}

type serviceImpl struct {
	dalService idal.IDAL
}

func NewService() IService {
	sv := &serviceImpl{}
	gormSession := db.GormFactory()
	sv.dalService = dal.New(
		dal.WithDB(gormSession),
	)
	return sv
}
