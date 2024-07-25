package domain

import (
	"database/sql"
	"reflect"

	"github.com/ericlagergren/decimal"
	"github.com/samwang0723/jarvis/internal/helper"
)

type Selection struct {
	StockID         string
	Name            string
	Category        string
	ExchangeDate    string
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

func ConvertSelectionList(sel any) []*Selection {
	var result []*Selection

	slice := reflect.ValueOf(sel)
	if slice.Kind() != reflect.Slice {
		panic("unsupported type")
	}

	for i := 0; i < slice.Len(); i++ {
		s := slice.Index(i).Interface()
		obj := mapToSelection(s)
		result = append(result, obj)
	}

	return result
}

func mapToSelection(s any) *Selection {
	var obj Selection

	// Manually handle the conversion of specific fields
	val := reflect.ValueOf(s).Elem()
	obj.Name = val.FieldByName("Name").Interface().(sql.NullString).String
	obj.StockID = val.FieldByName("StockID").String()
	obj.Category = val.FieldByName("Category").String()
	obj.ExchangeDate = val.FieldByName("ExchangeDate").String()
	obj.Volume = int(val.FieldByName("Volume").Float())
	obj.Trust = int(val.FieldByName("Trust").Float())
	obj.Foreign = int(val.FieldByName("Foreignc").Float())
	obj.Hedging = int(val.FieldByName("Hedging").Float())
	obj.Dealer = int(val.FieldByName("Dealer").Float())

	// Handle decimal.Big fields
	obj.Open = helper.DecimalToFloat32(
		val.FieldByName("Open").Interface().(decimal.Big),
	)
	obj.High = helper.DecimalToFloat32(
		val.FieldByName("High").Interface().(decimal.Big),
	)
	obj.Low = helper.DecimalToFloat32(
		val.FieldByName("Low").Interface().(decimal.Big),
	)
	obj.Close = helper.DecimalToFloat32(
		val.FieldByName("Close").Interface().(decimal.Big),
	)
	obj.PriceDiff = helper.DecimalToFloat32(
		val.FieldByName("PriceDiff").Interface().(decimal.Big),
	)
	obj.Concentration1 = helper.DecimalToFloat32(
		val.FieldByName("Concentration1").Interface().(decimal.Big),
	)
	obj.Concentration5 = helper.DecimalToFloat32(
		val.FieldByName("Concentration5").Interface().(decimal.Big),
	)
	obj.Concentration10 = helper.DecimalToFloat32(
		val.FieldByName("Concentration10").Interface().(decimal.Big),
	)
	obj.Concentration20 = helper.DecimalToFloat32(
		val.FieldByName("Concentration20").Interface().(decimal.Big),
	)
	obj.Concentration60 = helper.DecimalToFloat32(
		val.FieldByName("Concentration60").Interface().(decimal.Big),
	)

	return &obj
}
