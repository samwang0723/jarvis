package dto

type ThreePrimary struct {
	Model
	StockID         string `json:"StockID"`
	Date            string `json:"Date"`
	ForeignBuy      uint64 `json:"ForeignBuy"`
	ForeignSell     uint64 `json:"ForeignSell"`
	InvestTrustBuy  uint64 `json:"InvestTrustBuy"`
	InvestTrustSell uint64 `json:"InvestTrustSell"`
	DealerBuy       uint64 `json:"DealerBuy"`
	DealerSell      uint64 `json:"DealerSell"`
	HedgingBuy      uint64 `json:"HedgingBuy"`
	HedgingSell     uint64 `json:"HedgingSell"`
}
