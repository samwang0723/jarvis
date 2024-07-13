package domain

type Stock struct {
	Time
	ID       string
	Name     string
	Country  string
	Category string
	Market   string
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
