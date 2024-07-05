package domain

type ThreePrimary struct {
	ID
	StockID            string
	ExchangeDate       string
	ForeignTradeShares *int64
	TrustTradeShares   *int64
	DealerTradeShares  *int64
	HedgingTradeShares *int64
	Time
}

type ListThreePrimaryParams struct {
	Limit     int32
	Offset    int32
	StockID   string
	StartDate string
	EndDate   string
}
