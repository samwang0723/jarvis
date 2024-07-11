package domain

type Analysis struct {
	MA8       float32
	MA21      float32
	MA55      float32
	LastClose float32
	MV5       uint64
	MV13      uint64
	MV34      uint64
	Foreign   int64
	Trust     int64
	Hedging   int64
	Dealer    int64
}
