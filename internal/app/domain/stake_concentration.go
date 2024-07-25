package domain

type CalculationBase struct {
	Date        string
	TradeShares uint64
	Diff        int
}

type StakeConcentration struct {
	Time
	StockID         string  `json:"stockId"`
	Date            string  `json:"exchangeDate"`
	Diff            []int32 `json:"diff"`
	SumBuyShares    uint64  `json:"sumBuyShares"`
	SumSellShares   uint64  `json:"sumSellShares"`
	Concentration1  float32
	Concentration5  float32
	Concentration10 float32
	Concentration20 float32
	Concentration60 float32
	AvgBuyPrice     float32 `json:"avgBuyPrice"`
	AvgSellPrice    float32 `json:"avgSellPrice"`
	ID
}
