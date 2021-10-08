package dal

import (
	"context"
	"samwang0723/jarvis/entity"
)

func (i *dalImpl) CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error {
	err := i.db.Create(obj).Error
	return err
}

func (i *dalImpl) BatchCreateDailyClose(ctx context.Context, objs []*entity.DailyClose) error {
	err := i.db.Create(&objs).Error
	return err
}
