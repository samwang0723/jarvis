package domain

import (
	"database/sql"
	"reflect"

	"github.com/ericlagergren/decimal"
	"github.com/samwang0723/jarvis/internal/helper"
)

type StockPrice struct {
	ExchangeDate string  `json:"exchangeDate"`
	StockID      string  `json:"stockId"`
	Price        float32 `json:"price"`
}

type DailyClose struct {
	Time
	StockID      string
	ExchangeDate string  `json:"date"`
	TradedShares int64   `json:"tradeShares"`
	Transactions int64   `json:"transactions"`
	Turnover     int64   `json:"turnover"`
	Open         float32 `json:"open"`
	Close        float32 `json:"close"`
	High         float32 `json:"high"`
	Low          float32 `json:"low"`
	PriceDiff    float32 `json:"priceDiff"`
	ID
}

type ListDailyCloseParams struct {
	StartDate string
	StockID   string
	EndDate   string
	Limit     int32
	Offset    int32
}

func ConvertDailyCloseList(sel any) []*DailyClose {
	var result []*DailyClose

	slice := reflect.ValueOf(sel)
	if slice.Kind() != reflect.Slice {
		panic("unsupported type")
	}

	for i := 0; i < slice.Len(); i++ {
		s := slice.Index(i).Interface()
		obj := mapToDailyClose(s)
		result = append(result, obj)
	}

	return result
}

func mapToDailyClose(s any) *DailyClose {
	var obj DailyClose

	// Manually handle the conversion of specific fields
	val := reflect.ValueOf(s).Elem()
	obj.StockID = val.FieldByName("StockID").String()
	obj.ExchangeDate = val.FieldByName("ExchangeDate").String()
	obj.Close = helper.DecimalToFloat32(val.FieldByName("Close").Interface().(decimal.Big))
	obj.TradedShares = val.FieldByName("TradeShares").Interface().(sql.NullInt64).Int64

	return &obj
}
