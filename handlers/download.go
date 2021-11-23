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

// download job to run in workerpool
type downloadJob struct {
	ctx       context.Context
	date      string
	respChan  chan *[]interface{}
	rateLimit int
	origin    parser.Source
}

func (h *handlerImpl) CronDownload(ctx context.Context) error {
	return h.dataService.AddJob(ctx, "00 18 * * 1-5", func() {
		h.BatchingDownload(ctx, &dto.DownloadRequest{
			RewindLimit: 1,
			RateLimit:   2000,
		})
	})
}

// batching download all the historical stock data
func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	respChan := make(chan *[]interface{})

	go generateJob(ctx, req.RewindLimit*-1, parser.TwseDailyClose, req.RateLimit, respChan)
	go generateJob(ctx, req.RewindLimit*-1, parser.TpexDailyClose, req.RateLimit, respChan)

	go func() {
		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(4 * time.Hour):
				log.Warn("(BatchingDownload): timeout")
				return
			case <-ctx.Done():
				log.Warn("(BatchingDownload): context cancel")
				return
			case obj, ok := <-respChan:
				if ok {
					h.dataService.BatchUpsertDailyClose(ctx, obj)
				}
			}
		}
	}()
}

func generateJob(ctx context.Context, start int, origin parser.Source, rateLimit int, respChan chan *[]interface{}) {
	for i := start; i <= 0; i++ {
		var date string
		switch origin {
		case parser.TwseDailyClose:
			date = helper.ConvertDateStr(0, 0, i, helper.TwseDateFormat)
		case parser.TpexDailyClose:
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
				log.Debug("(BatchingDownload): generateJob goroutine exit!")
				return
			}
		}
	}
	log.Debug("(BatchingDownload): all download jobs sent!")
}

func (job *downloadJob) Do() error {
	c := crawler.New()
	p := parser.New()

	switch job.origin {
	case parser.TwseDailyClose:
		c.SetURL(icrawler.TwseDailyClose, job.date, icrawler.StockOnly)
	case parser.TpexDailyClose:
		c.SetURL(icrawler.TpexDailyClose, job.date)
	default:
		return fmt.Errorf("no recognized job source being specified: %s", job.origin)
	}
	config := parser.Config{
		ParseDay: &job.date,
		Capacity: 17,
		Type:     job.origin,
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
		log.Warn("(BatchingDownload): downloadJob - context cancelled!")
	}

	return nil
}
