// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package businessmodel

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/jarvis/internal/helper"
)

//nolint:nolintlint,gochecknoglobals
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// https://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_2330.tw
type Realtime struct {
	StockID   string  `json:"stockID"`
	Name      string  `json:"name"`
	Date      string  `json:"date"`
	ParseTime string  `json:"parseTime"`
	Open      float32 `json:"open"`
	Close     float32 `json:"close"`
	High      float32 `json:"high"`
	Low       float32 `json:"low"`
	Volume    uint64  `json:"volume"`
}

type rawData struct {
	MessageAry []rawBody `json:"msgArray"`
}

type rawBody struct {
	Buy5    string `json:"a"`
	Sell5   string `json:"b"`
	StockID string `json:"c"`
	Date    string `json:"d"`
	High    string `json:"h"`
	Low     string `json:"l"`
	Open    string `json:"o"`
	Close   string `json:"z"`
	Volume  string `json:"v"`
	Time    string `json:"t"`
	Name    string `json:"n"`
}

func (r *Realtime) UnmarshalJSON(data []byte) error {
	var raw rawData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.StockID = raw.MessageAry[0].StockID
	r.Date = raw.MessageAry[0].Date
	//nolint:nolintlint,errcheck
	r.Open, _ = helper.StringToFloat32(raw.MessageAry[0].Open)
	//nolint:nolintlint,errcheck
	r.Close, _ = helper.StringToFloat32(raw.MessageAry[0].Close)
	//nolint:nolintlint,errcheck
	r.High, _ = helper.StringToFloat32(raw.MessageAry[0].High)
	//nolint:nolintlint,errcheck
	r.Low, _ = helper.StringToFloat32(raw.MessageAry[0].Low)
	//nolint:nolintlint,errcheck
	r.Volume, _ = helper.StringToUint64(raw.MessageAry[0].Volume)
	r.ParseTime = raw.MessageAry[0].Time

	return nil
}
