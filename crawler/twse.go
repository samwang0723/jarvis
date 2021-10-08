package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"samwang0723/jarvis/helper"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

const (
	All           = "ALL"
	StockOnly     = "ALLBUT0999"
	DailyClose    = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=%s"
	ThreePrimary  = "http://www.tse.com.tw/fund/T86?response=csv&date=%s&selectType=%s"
	OperatingDays = "https://www.twse.com.tw/holidaySchedule/holidaySchedule?response=csv&queryYear=%d"
)

type TwseStock struct {
	url string
}

func (twse *TwseStock) Fetch() (io.Reader, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", twse.url, nil)
	if err != nil {
		return nil, fmt.Errorf("TWSE new fetch request initialize error: %v\n", err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TWSE fetch request error: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TWSE fetch status error: %v\n", resp.StatusCode)
	}

	csvfile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TWSE fetch unable to read body: %v\n", err)
	}

	log.Printf("TWSE download completed (%s), URL: %s, Header: %v\n",
		helper.ReadableSize(len(csvfile), 2), twse.url, resp.Header)
	raw := bytes.NewBuffer(csvfile)
	reader := transform.NewReader(raw, traditionalchinese.Big5.NewDecoder())

	return reader, nil
}

func (twse *TwseStock) SetURL(template string, date string, queryType string) {
	twse.url = fmt.Sprintf(template, date, queryType)
}
