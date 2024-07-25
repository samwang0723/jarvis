package domain

type Stock struct {
	Time
	ID       string `json:"stockId"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Category string `json:"category"`
	Market   string `json:"market"`
}

type ListStocksParams struct {
	Country         string
	Name            string
	Category        string
	StockIDs        []string
	Limit           int32
	Offset          int32
	FilterByStockID bool
}
