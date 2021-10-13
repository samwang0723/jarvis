package handlers

import (
	"context"
	"log"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/helper"
	"samwang0723/jarvis/parser"
	"sync"
	"time"
)

func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	wg := &sync.WaitGroup{}
	respChan := make(chan *[]interface{})
	wg.Add(2)
	go twse(ctx, wg, respChan, req)
	go tpex(ctx, wg, respChan, req)
	go func() {
		wg.Wait()
		close(respChan)
	}()
	for {
		if obj, ok := <-respChan; ok {
			h.dataService.BatchCreateDailyClose(ctx, obj)
		} else {
			break
		}
	}
}

func twse(ctx context.Context, wg *sync.WaitGroup, respChan chan *[]interface{}, req *dto.DownloadRequest) {
	c := crawler.New()
	p := parser.New()
	defer wg.Done()
	for i := req.RewindLimit; i <= 0; i++ {
		dayString := helper.ConvertDateStr(0, 0, i, helper.TwseDateFormat)
		if len(dayString) == 0 {
			continue
		}
		c.SetURL(icrawler.TwseDailyClose, dayString, icrawler.StockOnly)
		stream, err := c.Fetch()
		if err != nil {
			log.Printf("twse DailyClose fetch error: %s\n", err)
			continue
		}
		err = p.Parse(parser.Config{
			ParseDay: &dayString,
			Capacity: 17,
			Type:     parser.TwseDailyClose,
		}, stream)
		if err != nil {
			continue
		}
		respChan <- p.Flush()

		time.Sleep(time.Duration(req.RateLimit) * time.Millisecond)
	}
}

func tpex(ctx context.Context, wg *sync.WaitGroup, respChan chan *[]interface{}, req *dto.DownloadRequest) {
	c := crawler.New()
	p := parser.New()
	defer wg.Done()
	for i := req.RewindLimit; i <= 0; i++ {
		dayString := helper.ConvertDateStr(0, 0, i, helper.TpexDateFormat)
		if len(dayString) == 0 {
			continue
		}
		c.SetURL(icrawler.TpexDailyClose, dayString)
		stream, err := c.Fetch()
		if err != nil {
			continue
		}
		err = p.Parse(parser.Config{
			ParseDay: &dayString,
			Capacity: 19,
			Type:     parser.TpexDailyClose,
		}, stream)
		if err != nil {
			continue
		}
		respChan <- p.Flush()

		time.Sleep(time.Duration(req.RateLimit) * time.Millisecond)
	}
}
