package dal

import (
	"context"
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/entity"
	"strings"

	"gorm.io/gorm"
)

func (i *dalImpl) CreateDailyClose(ctx context.Context, obj *entity.DailyClose) error {
	err := i.db.Create(obj).Error
	return err
}

func (i *dalImpl) BatchCreateDailyClose(ctx context.Context, objs []*entity.DailyClose) error {
	err := i.db.Create(&objs).Error
	return err
}

func (i *dalImpl) ListDailyClose(ctx context.Context, offset int, limit int,
	searchParams *idal.ListDailyCloseSearchParams) (objs []*entity.DailyClose, totalCount int64, err error) {
	query := i.db.Model(&entity.DailyClose{})
	query = buildQueryFromListDailyCloseSearchParams(query, searchParams)
	err = query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&objs).Error
	if err != nil {
		return nil, 0, err
	}
	return objs, totalCount, nil
}

func buildQueryFromListDailyCloseSearchParams(query *gorm.DB, params *idal.ListDailyCloseSearchParams) *gorm.DB {
	if params == nil {
		return query
	}
	query = query.Where("exchange_date >= ?", params.Start)
	if params.End != nil {
		dateStr := *params.End
		query = query.Where("exchange_date < ?", dateStr)
	}
	if params.StockIDs != nil {
		query = query.Where("stock_id IN (" + strings.Join(*params.StockIDs, ",") + ")")
	}
	return query
}
