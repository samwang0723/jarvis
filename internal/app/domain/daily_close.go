package domain

type StockPrice struct {
	ExchangeDate string
	StockID      string
	Price        float32
}

type DailyClose struct {
	ID
	StockID      string
	Date         string
	TradedShares int64
	Transactions int64
	Turnover     int64
	Open         float32
	Close        float32
	High         float32
	Low          float32
	PriceDiff    float32
	Time
}

type ListDailyCloseParams struct {
	Limit     int32
	Offset    int32
	StartDate string
	StockID   string
	EndDate   string
}
