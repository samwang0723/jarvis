package handlers

import (
	"context"
	"io"
	"log"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/helper"
	"samwang0723/jarvis/parser"
	"time"
)

func (h *handlerImpl) BatchingDownload(ctx context.Context, rewindLimit int, rateLimit int) {
	for i := rewindLimit; i <= 0; i++ {
		d := helper.ConvertDateStr(0, 0, i)
		twse := crawler.New()
		twse.SetURL(icrawler.TwseDailyClose, d, icrawler.StockOnly)
		fetchStream, err := twse.Fetch()
		if err != nil {
			log.Printf("DailyClose fetch error: %s\n", err)
			continue
		}
		resp, err := parse(ctx, d, fetchStream)
		if err != nil || len(*resp) == 0 {
			log.Printf("DailyClose parse error: %v\n", err)
			continue
		}
		err = h.dataService.BatchCreateDailyClose(ctx, resp)
		if err != nil {
			log.Printf("DailyClose persistent storage failed: %v\n", err)
			continue
		}
		time.Sleep(time.Duration(rateLimit) * time.Millisecond)
	}
}

func parse(ctx context.Context, day string, stream io.Reader) (*[]interface{}, error) {
	p := parser.New()
	data := &[]interface{}{}
	p.SetDataSource(data)
	config := parser.Config{
		ParseDay: &day,
		Capacity: 17,
		Type:     parser.TwseDailyClose,
	}
	return p.Parse(config, stream)
}
