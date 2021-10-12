package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/helper"
	"strings"
)

func (p *parserImpl) Parse(config Config, in io.Reader) (*[]interface{}, error) {
	if p.result == nil {
		return nil, fmt.Errorf("didn't initialized the result map\n")
	}
	if config.ParseDay == nil {
		return nil, fmt.Errorf("parse day missing\n")
	}

	reader := csv.NewReader(in)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if len(record) == 0 {
			continue
		}
		idColumn := record[0]
		if helper.IsInteger(idColumn[0:2]) && config.Capacity == len(record) {
			switch config.Type {
			case TwseDailyClose:
				*p.result = append(*p.result, twseToEntity(*config.ParseDay, record))
			case TwseThreePrimary:
			case TpexDailyClose:
				*p.result = append(*p.result, tpexToEntity(*config.ParseDay, record))
			}
		}
	}
	return p.result, nil
}

func twseToEntity(day string, data []string) *entity.DailyClose {
	id := data[0]
	dailyClose := &entity.DailyClose{
		StockID:      id,
		Date:         day,
		TradedShares: helper.ToUint64(strings.Replace(data[2], ",", "", -1)),
		Transactions: helper.ToUint64(strings.Replace(data[3], ",", "", -1)),
		Turnover:     helper.ToUint64(strings.Replace(data[4], ",", "", -1)),
		Open:         helper.ToFloat32(data[5]),
		High:         helper.ToFloat32(data[6]),
		Low:          helper.ToFloat32(data[7]),
		Close:        helper.ToFloat32(data[8]),
		PriceDiff:    helper.ToFloat32(fmt.Sprintf("%s%s", data[9], data[10])),
	}
	return dailyClose
}

func tpexToEntity(day string, data []string) *entity.DailyClose {
	id := data[0]
	dailyClose := &entity.DailyClose{
		StockID:      id,
		Date:         day,
		TradedShares: helper.ToUint64(strings.Replace(data[7], ",", "", -1)),
		Transactions: helper.ToUint64(strings.Replace(data[9], ",", "", -1)),
		Turnover:     helper.ToUint64(strings.Replace(data[8], ",", "", -1)),
		Open:         helper.ToFloat32(data[4]),
		High:         helper.ToFloat32(data[5]),
		Low:          helper.ToFloat32(data[6]),
		Close:        helper.ToFloat32(data[2]),
		PriceDiff:    helper.ToFloat32(data[3]),
	}
	return dailyClose
}
