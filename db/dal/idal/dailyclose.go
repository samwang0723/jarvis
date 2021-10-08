package idal

import (
	"context"
	"samwang0723/jarvis/entity"
)

type IDailyCloseDAL interface {
	CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error
	BatchCreateDailyClose(ctx context.Context, objs []*entity.DailyClose) error
}
