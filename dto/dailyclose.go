package dto

type DailyClose struct {
	Model
	StockID      string  `json:"StockID"`
	Date         string  `json:"Date"`
	TradedShares uint64  `json:"TradedShares,omitempty"`
	Transactions uint64  `json:"Transactions,omitempty"`
	Turnover     uint64  `json:"Turnover,omitempty"`
	Open         float32 `json:"Open"`
	Close        float32 `json:"Close"`
	High         float32 `json:"High"`
	Low          float32 `json:"Low"`
	PriceDiff    float32 `json:"PriceDiff"`
}
