package handlers

import (
	"context"
	"fmt"
	"log"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/helper"
	"samwang0723/jarvis/parser"
	"time"
)

func (h *handlerImpl) BatchingDownload(ctx context.Context, rewindLimit int, rateLimit int) {
	for i := rewindLimit; i <= 0; i++ {
		d := helper.ConvertDateStr(0, 0, i)
		resp, err := downloadDailyCloses(ctx, d)
		if err != nil {
			log.Fatalf("download DailyClose error: %v\n", err)
		}
		err = h.dataService.BatchCreateDailyClose(ctx, resp)
		if err != nil {
			log.Fatalf("DailyClose persistent storage failed: %v\n", err)
		}
		time.Sleep(time.Duration(rateLimit) * time.Millisecond)
	}
}

func downloadDailyCloses(ctx context.Context, day string) (map[string]interface{}, error) {
	var twse crawler.ICrawler
	twse = &crawler.TwseStock{}
	twse.SetURL(crawler.DailyClose, day, crawler.StockOnly)
	fetchStream, err := twse.Fetch()
	if err != nil {
		return nil, fmt.Errorf("DailyClose fetch error: %s\n", err)
	}

	var p parser.IParser
	p = &parser.CsvSource{Tag: day}
	data := map[string]interface{}{}
	p.SetDataSource(data)
	config := parser.Config{
		StartInteger: true,
		Capacity:     17,
		Type:         parser.TwseDailyClose,
	}
	resp, err := p.Parse(config, fetchStream)
	if err != nil {
		return nil, fmt.Errorf("DailyClose parse error: %s\n", err)
	}

	return resp, nil
}
