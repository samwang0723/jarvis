package domain

type CalculationBase struct {
	Date        string
	TradeShares uint64
	Diff        int
}

type StakeConcentration struct {
	ID
	StockID         string
	Date            string
	SumBuyShares    uint64
	SumSellShares   uint64
	AvgSellPrice    float32
	Concentration1  float32
	Concentration5  float32
	Concentration10 float32
	Concentration20 float32
	Concentration60 float32
	AvgBuyPrice     float32
	Diff            []int32
	Time
}
