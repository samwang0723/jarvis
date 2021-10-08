package entity

type Stock struct {
	Model

	StockID string `gorm:"column:stock_id"`
	Name    string `gorm:"column:name"`
	Country string `gorm:"column:country"`
}

func (Stock) TableName() string {
	return "stocks"
}
