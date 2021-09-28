package handlers

import (
	"context"
	"fmt"
	"reflect"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/parser"
)

func DownloadDailyCloses(ctx context.Context, day string) (map[string]*dto.DailyClose, error) {
	var twse crawler.Crawler
	twse = &crawler.TwseStock{}
	twse.SetURL(crawler.DailyClose, day, crawler.StockOnly)
	io, err := twse.Fetch()
	if err != nil {
		return nil, fmt.Errorf("Fetch error: %s\n", err)
	}

	var handler parser.Parser
	handler = &parser.CsvHandler{Tag: day}
	data := map[string]interface{}{}
	handler.SetDataSource(data)
	config := parser.Config{
		StartInteger: true,
		Capacity:     17,
		Type:         parser.TwseDailyClose,
	}
	resp, err := handler.Parse(config, io)
	if err != nil {
		return nil, fmt.Errorf("Parse error: %s\n", err)
	}

	// Replicate the value from interface to *dto.DailyClose
	dailyCloses := map[string]*dto.DailyClose{}
	for k, v := range resp {
		if val, ok := v.(*dto.DailyClose); ok {
			dailyCloses[k] = val
		} else {
			return nil, fmt.Errorf("Cannot cast interface to *dto.DailyClose: %v\n", reflect.TypeOf(v).Elem())
		}
	}

	return dailyCloses, nil
}
