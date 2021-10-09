package idal

import (
	"context"
	"samwang0723/jarvis/entity"
)

type ListDailyCloseSearchParams struct {
	StockIDs *[]string
	Start    string
	End      *string
}

type IDailyCloseDAL interface {
	CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error
	BatchCreateDailyClose(ctx context.Context, objs []*entity.DailyClose) error
	ListDailyClose(ctx context.Context, offset int, limit int,
		searchParams *ListDailyCloseSearchParams) (objs []*entity.DailyClose, totalCount int64, err error)
}
