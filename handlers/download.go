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
package handlers

import (
	"context"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/helper"
	log "samwang0723/jarvis/logger"
	"samwang0723/jarvis/parser"

	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	wg := &sync.WaitGroup{}
	respChan := make(chan *[]interface{})
	wg.Add(2)
	go download(ctx, wg, respChan, req, twse)
	go download(ctx, wg, respChan, req, tpex)
	go func() {
		wg.Wait()
		close(respChan)
	}()
	for {
		select {
		case <-ctx.Done():
			log.Warn("terminate the downloading process")
			return
		case obj, ok := <-respChan:
			if ok {
				h.dataService.BatchCreateDailyClose(ctx, obj)
			} else {
				return
			}
		}
	}
}

func download(ctx context.Context,
	wg *sync.WaitGroup,
	respChan chan *[]interface{},
	req *dto.DownloadRequest,
	fn func(int, chan *[]interface{})) {

	defer wg.Done()

	index := req.RewindLimit * -1
	for {
		fn(index, respChan)
		// calculate count
		index++
		if index > 0 {
			return
		}

		// rate limit protection and context.cancel
		select {
		case <-time.After(time.Duration(req.RateLimit) * time.Millisecond):
			break
		case <-ctx.Done():
			log.Warn("download: context cancelled!")
			return
		}
	}
}

func twse(index int, respChan chan *[]interface{}) {
	c := crawler.New()
	p := parser.New()
	dayString := helper.ConvertDateStr(0, 0, index, helper.TwseDateFormat)
	if len(dayString) == 0 {
		return
	}
	c.SetURL(icrawler.TwseDailyClose, dayString, icrawler.StockOnly)
	stream, err := c.Fetch()
	if err != nil {
		sentry.CaptureException(err)
		log.Errorf("twse DailyClose fetch error: %s\n", err)
		return
	}
	err = p.Parse(parser.Config{
		ParseDay: &dayString,
		Capacity: 17,
		Type:     parser.TwseDailyClose,
	}, stream)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	respChan <- p.Flush()
}

func tpex(index int, respChan chan *[]interface{}) {
	c := crawler.New()
	p := parser.New()
	dayString := helper.ConvertDateStr(0, 0, index, helper.TpexDateFormat)
	if len(dayString) == 0 {
		return
	}
	c.SetURL(icrawler.TpexDailyClose, dayString)
	stream, err := c.Fetch()
	if err != nil {
		sentry.CaptureException(err)
		log.Errorf("tpex DailyClose fetch error: %s\n", err)
		return
	}
	err = p.Parse(parser.Config{
		ParseDay: &dayString,
		Capacity: 19,
		Type:     parser.TpexDailyClose,
	}, stream)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	respChan <- p.Flush()
}
