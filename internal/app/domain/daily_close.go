package domain

type DailyClose struct {
	ID
	StockID      string
	Date         string
	TradedShares uint64
	Transactions uint64
	Turnover     uint64
	Open         float32
	Close        float32
	High         float32
	Low          float32
	PriceDiff    float32
	Time
}
