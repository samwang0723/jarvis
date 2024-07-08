package domain

type Selection struct {
	StockID         string
	Name            string
	Category        string
	Date            string
	Open            float32
	High            float32
	Low             float32
	Close           float32
	PriceDiff       float32
	Concentration1  float32
	Concentration5  float32
	Concentration10 float32
	Concentration20 float32
	Concentration60 float32
	Volume          int
	Trust           int
	Foreign         int
	Hedging         int
	Dealer          int
	Trust10         int
	Foreign10       int
	QuoteChange     float32
}
