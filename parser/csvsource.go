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
		if helper.IsInteger(record[0]) && config.Capacity == len(record) {
			switch config.Type {
			case TwseDailyClose:
				*p.result = append(*p.result, twseToEntity(*config.ParseDay, record))
			case TwseThreePrimary:
			}
		}
	}
	return p.result, nil
}

func (p *parserImpl) SetDataSource(source *[]interface{}) {
	p.result = source
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
