package domain

type CalculationBase struct {
	Date        string
	TradeShares uint64
	Diff        int
}

type StakeConcentration struct {
	Time
	StockID         string
	Date            string
	Diff            []int32
	SumBuyShares    uint64
	SumSellShares   uint64
	Concentration1  float32
	Concentration5  float32
	Concentration10 float32
	Concentration20 float32
	Concentration60 float32
	AvgBuyPrice     float32
	AvgSellPrice    float32
	ID
}
