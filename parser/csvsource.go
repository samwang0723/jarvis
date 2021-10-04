package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/helper"
	"strings"
)

type CsvSource struct {
	Tag    string
	result map[string]interface{}
}

func (handler *CsvSource) Parse(config Config, in io.Reader) (map[string]interface{}, error) {
	if handler.result == nil {
		return nil, fmt.Errorf("Didn't initialized the result map\n")
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
		if config.StartInteger && helper.IsInteger(record[0]) && config.Capacity == len(record) {
			switch config.Type {
			case TwseDailyClose:
				handler.storeTwseDailyClose(record)
			case TwseThreePrimary:
			}
		}
	}
	return handler.result, nil
}

func (handler *CsvSource) SetDataSource(source map[string]interface{}) {
	handler.result = source
}

func (handler *CsvSource) storeTwseDailyClose(data []string) {
	id := data[0]
	dailyclose := &dto.DailyClose{
		StockID:      id,
		Date:         handler.Tag,
		TradedShares: helper.ToUint64(strings.Replace(data[2], ",", "", -1)),
		Transactions: helper.ToUint64(strings.Replace(data[3], ",", "", -1)),
		Turnover:     helper.ToUint64(strings.Replace(data[4], ",", "", -1)),
		Open:         helper.ToFloat32(data[5]),
		High:         helper.ToFloat32(data[6]),
		Low:          helper.ToFloat32(data[7]),
		Close:        helper.ToFloat32(data[8]),
		PriceDiff:    helper.ToFloat32(fmt.Sprintf("%s%s", data[9], data[10])),
	}
	handler.result[id] = dailyclose
}
