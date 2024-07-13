package domain

import (
	"database/sql"
	"reflect"

	"github.com/ericlagergren/decimal"
	"github.com/mitchellh/mapstructure"
	"github.com/samwang0723/jarvis/internal/helper"
)

type StockPrice struct {
	ExchangeDate string
	StockID      string
	Price        float32
}

type DailyClose struct {
	Time
	StockID      string
	ExchangeDate string
	TradedShares int64
	Transactions int64
	Turnover     int64
	Open         float32
	Close        float32
	High         float32
	Low          float32
	PriceDiff    float32
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
		obj, err := mapToDailyClose(s)
		if err != nil {
			panic(err)
		}
		result = append(result, obj)
	}

	return result
}

func mapToDailyClose(s any) (*DailyClose, error) {
	var obj DailyClose
	err := mapstructure.Decode(s, &obj)
	if err != nil {
		return nil, err
	}

	// Manually handle the conversion of specific fields
	val := reflect.ValueOf(s).Elem()
	obj.StockID = val.FieldByName("StockID").String()
	obj.ExchangeDate = val.FieldByName("ExchangeDate").String()
	obj.Close = helper.DecimalToFloat32(val.FieldByName("Close").Interface().(decimal.Big))
	obj.TradedShares = val.FieldByName("TradeShares").Interface().(sql.NullInt64).Int64

	return &obj, nil
}
