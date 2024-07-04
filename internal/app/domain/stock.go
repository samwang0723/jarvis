package domain

type Stock struct {
	ID       string
	Name     string
	Country  string
	Category string
	Market   string
	Time
}

type ListStocksParams struct {
	Limit           int32
	Offset          int32
	Country         string
	StockIDs        []string
	FilterByStockID bool
	Name            string
	Category        string
}
