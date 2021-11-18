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
	"fmt"
	"samwang0723/jarvis/concurrent"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/helper"
	log "samwang0723/jarvis/logger"
	"samwang0723/jarvis/parser"
	"time"

	"github.com/getsentry/sentry-go"
)

type source int

//go:generate stringer -type=source
const (
	TwseDailyClose source = iota
	TwseThreePrimary
	TpexDailyClose
)

// download job to run in workerpool
type downloadJob struct {
	ctx       context.Context
	date      string
	respChan  chan *[]interface{}
	rateLimit int
	origin    source
}

// batching download all the historical stock data
func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	respChan := make(chan *[]interface{})
	defer close(respChan)

	go generateJob(ctx, req.RewindLimit*-1, TwseDailyClose, req.RateLimit, respChan)
	go generateJob(ctx, req.RewindLimit*-1, TpexDailyClose, req.RateLimit, respChan)

	for {
		select {
		case <-ctx.Done():
			log.Warn("=== terminate the downloading process ===")
			return
		case obj, ok := <-respChan:
			if ok {
				h.dataService.BatchCreateDailyClose(ctx, obj)
			}
		}
	}
}

func generateJob(ctx context.Context, start int, origin source, rateLimit int, respChan chan *[]interface{}) {
	for i := start; i < 0; i++ {
		var date string
		switch origin {
		case TwseDailyClose:
			date = helper.ConvertDateStr(0, 0, i, helper.TwseDateFormat)
		case TpexDailyClose:
			date = helper.ConvertDateStr(0, 0, i, helper.TpexDateFormat)
		}

		if len(date) > 0 {
			job := &downloadJob{
				ctx:       ctx,
				date:      date,
				respChan:  respChan,
				rateLimit: rateLimit,
				origin:    origin,
			}
			select {
			case concurrent.JobQueue <- job:
			case <-ctx.Done():
				log.Debug("download: generateJob goroutine exit!")
				return
			}
		}
	}
}

func (job *downloadJob) Do() error {
	c := crawler.New()
	p := parser.New()
	var config parser.Config

	switch job.origin {
	case TwseDailyClose:
		c.SetURL(icrawler.TwseDailyClose, job.date, icrawler.StockOnly)
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 17,
			Type:     parser.TwseDailyClose,
		}
	case TpexDailyClose:
		c.SetURL(icrawler.TpexDailyClose, job.date)
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 19,
			Type:     parser.TpexDailyClose,
		}
	default:
		return fmt.Errorf("no recognized job source being specified: %s", job.origin)
	}

	stream, err := c.Fetch(job.ctx)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("(%s/%s) fetch error: %+v", job.origin, job.date, err)
	}
	err = p.Parse(config, stream)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("(%s/%s) parse failed, err: %+v", job.origin, job.date, err)
	}

	job.respChan <- p.Flush()

	// rate limit protection and context.cancel
	select {
	case <-time.After(time.Duration(job.rateLimit) * time.Millisecond):
	case <-job.ctx.Done():
		log.Warn("download: context cancelled!")
	}

	return nil
}
