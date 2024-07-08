package domain

import (
	"errors"

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

var errEmptyData = errors.New("cannot unmarshal empty data")

func (r *Realtime) UnmarshalJSON(data []byte) error {
	var raw rawData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw.MessageAry) == 0 {
		return errEmptyData
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
