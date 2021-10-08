package entity

type DailyClose struct {
	Model

	StockID      string  `gorm:"column:stock_id"`
	Date         string  `gorm:"column:exchange_date"`
	TradedShares uint64  `gorm:"column:trade_shares"`
	Transactions uint64  `gorm:"column:transactions"`
	Turnover     uint64  `gorm:"column:turnover"`
	Open         float32 `gorm:"column:open"`
	Close        float32 `gorm:"column:close"`
	High         float32 `gorm:"column:high"`
	Low          float32 `gorm:"column:low"`
	PriceDiff    float32 `gorm:"column:price_diff"`
}

func (DailyClose) TableName() string {
	return "daily_closes"
}
