package dto

import "samwang0723/jarvis/entity"

type ListDailyCloseSearchParams struct {
	StockIDs *[]string `json:"StockIDs,omitempty"`
	Start    string    `json:"Start"`
	End      *string   `json:"End,omitempty"`
}

type ListDailyCloseRequest struct {
	Offset       int                         `json:"Offset"`
	Limit        int                         `json:"Limit"`
	SearchParams *ListDailyCloseSearchParams `json:"SearchParams"`
}

type ListDailyCloseResponse struct {
	Offset     int                  `json:"Offset"`
	Limit      int                  `json:"Limit"`
	TotalCount int                  `json:"TotalCount"`
	Entries    []*entity.DailyClose `json:"Entries"`
}
