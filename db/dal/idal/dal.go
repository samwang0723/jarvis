package idal

const (
	MaxRow = 1000
)

type IDAL interface {
	IStockDAL
	IDailyCloseDAL
}
