package entity

type ThreePrimary struct {
	Model

	StockID         string `gorm:"column:stock_id"`
	Date            string `gorm:"column:exchange_date"`
	ForeignBuy      uint64 `gorm:"column:foreign_buy"`
	ForeignSell     uint64 `gorm:"column:foreign_sell"`
	InvestTrustBuy  uint64 `gorm:"column:invest_trust_buy"`
	InvestTrustSell uint64 `gorm:"column:invest_trust_sell"`
	DealerBuy       uint64 `gorm:"column:dealer_buy"`
	DealerSell      uint64 `gorm:"column:dealer_sell"`
	HedgingBuy      uint64 `gorm:"column:hedging_buy"`
	HedgingSell     uint64 `gorm:"column:hedging_sell"`
}
