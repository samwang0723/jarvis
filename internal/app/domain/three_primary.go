package domain

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type ThreePrimary struct {
	Time
	StockID            string `json:"stockId"`
	ExchangeDate       string `json:"exchangeDate"`
	ForeignTradeShares int64  `json:"foreignTradeShares"`
	TrustTradeShares   int64  `json:"trustTradeShares"`
	DealerTradeShares  int64  `json:"dealerTradeShares"`
	HedgingTradeShares int64  `json:"hedgingTradeShares"`
	ID
}

type ListThreePrimaryParams struct {
	StockID   string
	StartDate string
	EndDate   string
	Limit     int32
	Offset    int32
}

func ConvertThreePrimaryList(sel any) []*ThreePrimary {
	var result []*ThreePrimary

	slice := reflect.ValueOf(sel)
	if slice.Kind() != reflect.Slice {
		panic("unsupported type")
	}

	for i := 0; i < slice.Len(); i++ {
		s := slice.Index(i).Interface()
		obj, err := mapToThreePrimary(s)
		if err != nil {
			panic(err)
		}
		result = append(result, obj)
	}

	return result
}

func mapToThreePrimary(s any) (*ThreePrimary, error) {
	var obj ThreePrimary
	err := mapstructure.Decode(s, &obj)
	if err != nil {
		return nil, err
	}

	// Manually handle the conversion of specific fields
	val := reflect.ValueOf(s).Elem()
	obj.StockID = val.FieldByName("StockID").String()
	obj.ExchangeDate = val.FieldByName("ExchangeDate").String()
	obj.ForeignTradeShares = int64(val.FieldByName("ForeignTradeShares").Float())
	obj.TrustTradeShares = int64(val.FieldByName("TrustTradeShares").Float())
	obj.DealerTradeShares = int64(val.FieldByName("DealerTradeShares").Float())
	obj.HedgingTradeShares = int64(val.FieldByName("HedgingTradeShares").Float())

	return &obj, nil
}
