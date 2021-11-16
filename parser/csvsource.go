// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/helper"
	"strings"
)

func (p *parserImpl) Parse(config Config, in io.Reader) error {
	if p.result == nil {
		return fmt.Errorf("didn't initialized the result map\n")
	}
	if config.ParseDay == nil {
		return fmt.Errorf("parse day missing\n")
	}

	originLen := len(*p.result)
	updatedLen := originLen

	reader := csv.NewReader(in)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	//override to standarize date string
	date := helper.UnifiedDateStr(*config.ParseDay)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if len(record) > 0 && len(record[0]) < 6 {
			idColumn := record[0]
			if helper.IsInteger(idColumn[0:2]) && config.Capacity == len(record) {
				switch config.Type {
				case TwseDailyClose:
					*p.result = append(*p.result, twseToEntity(date, record))
					updatedLen++
				case TwseThreePrimary:
				case TpexDailyClose:
					*p.result = append(*p.result, tpexToEntity(date, record))
					updatedLen++
				}
			}
		}
	}

	if updatedLen <= originLen {
		return fmt.Errorf("empty parsing results\n")
	}

	return nil
}

func (p *parserImpl) Flush() *[]interface{} {
	res := *p.result
	p.result = &[]interface{}{}
	return &res
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
