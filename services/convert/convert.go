package convert

import (
	"samwang0723/jarvis/db/dal/idal"
	"samwang0723/jarvis/dto"
)

func ListDailyCloseSearchParamsDTOToDAL(obj *dto.ListDailyCloseSearchParams) *idal.ListDailyCloseSearchParams {
	res := &idal.ListDailyCloseSearchParams{
		Start: obj.Start,
	}
	if obj.StockIDs != nil {
		res.StockIDs = obj.StockIDs
	}
	if obj.End != nil {
		res.End = obj.End
	}
	return res
}
