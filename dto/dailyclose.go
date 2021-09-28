package dto

import "time"

type DailyClose struct {
	StockID      string     `json:"StockID"`
	Date         string     `json:"Date"`
	TradedShares uint64     `json:"TradedShares"`
	Transactions uint64     `json:"Transactions"`
	Turnover     uint64     `json:"Turnover"`
	Open         float32    `json:"Open"`
	Close        float32    `json:"Close"`
	High         float32    `json:"High"`
	Low          float32    `json:"Low"`
	PriceDiff    float32    `json:"PriceDiff"`
	CreatedAt    *time.Time `json:"CreatedAt"`
	UpdatedAt    *time.Time `json:"UpdatedAt"`
}
