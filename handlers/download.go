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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samwang0723/jarvis/concurrent"
	"github.com/samwang0723/jarvis/crawler"
	"github.com/samwang0723/jarvis/crawler/icrawler"
	"github.com/samwang0723/jarvis/crawler/proxy"
	"github.com/samwang0723/jarvis/dto"
	"github.com/samwang0723/jarvis/entity"
	"github.com/samwang0723/jarvis/helper"
	log "github.com/samwang0723/jarvis/logger"
	"github.com/samwang0723/jarvis/parser"

	"github.com/getsentry/sentry-go"
)

const (
	StartCronjob = "START_CRON"
)

// download job to run in workerpool
type downloadJob struct {
	ctx       context.Context
	date      string
	stockId   string
	respChan  chan *[]interface{}
	rateLimit int32
	origin    parser.Source
}

func (h *handlerImpl) CronDownload(ctx context.Context, req *dto.StartCronjobRequest) (*dto.StartCronjobResponse, error) {
	envCron := os.Getenv(StartCronjob)
	startCron, err := strconv.ParseBool(envCron)
	if err != nil || !startCron {
		return &dto.StartCronjobResponse{
			Code:     401,
			Error:    "Unauthorized",
			Messages: "Environment not allowed to trigger Cronjob",
		}, err
	}

	// create a separate context since it's not rely on parent grpc.Dial()
	longLiveCtx := context.Background()
	err = h.dataService.AddJob(longLiveCtx, req.Schedule, func() {
		h.BatchingDownload(longLiveCtx, &dto.DownloadRequest{
			RewindLimit: 0,
			RateLimit:   3000,
			Types:       req.Types,
		})
	})

	if err != nil {
		return &dto.StartCronjobResponse{
			Code:     400,
			Error:    "Bad Request",
			Messages: fmt.Sprintf("Failed to start the schedule: %s with Types: %+v", req.Schedule, req.Types),
		}, err
	}

	return &dto.StartCronjobResponse{
		Code:     200,
		Messages: fmt.Sprintf("Successfully started the schedule: %s with Types: %+v", req.Schedule, req.Types),
	}, nil
}

func (h *handlerImpl) StockListDownload(ctx context.Context) {
	respChan := make(chan *[]interface{})
	types := []parser.Source{
		parser.TwseStockList,
		parser.TpexStockList,
	}
	go func() {
		for _, t := range types {
			job := &downloadJob{
				ctx:       ctx,
				respChan:  respChan,
				rateLimit: 1000,
				origin:    t,
			}
			select {
			case concurrent.JobQueue <- job:
			case <-ctx.Done():
				log.Debug("(StockListDownload): generateJob goroutine exit!")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(20 * time.Minute):
				log.Warn("(StockListDownload): timeout")
				return
			case <-ctx.Done():
				log.Warn("(StockListDownload): context cancel")
				return
			case objs, ok := <-respChan:
				if ok {
					h.dataService.BatchUpsertStocks(ctx, objs)
				}
			}
		}
	}()
}

// batching download all the historical stock data
func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	dailyCloseChan := make(chan *[]interface{})
	stakeConcentrationChan := make(chan *[]interface{})
	threePrimaryChan := make(chan *[]interface{})

	for _, t := range req.Types {
		switch t {
		case dto.DailyClose:
			go h.generateJob(ctx, parser.TwseDailyClose, req, dailyCloseChan)
			go h.generateJob(ctx, parser.TpexDailyClose, req, dailyCloseChan)
		case dto.ThreePrimary:
			go h.generateJob(ctx, parser.TwseThreePrimary, req, threePrimaryChan)
			go h.generateJob(ctx, parser.TpexThreePrimary, req, threePrimaryChan)
		case dto.Concentration:
			go h.generateJob(ctx, parser.StakeConcentration, req, stakeConcentrationChan)
		}
	}

	go func() {
		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(8 * time.Hour):
				log.Warn("(BatchingDownload): timeout")
				return
			case <-ctx.Done():
				log.Warn("(BatchingDownload): context cancel")
				return
			case objs, ok := <-dailyCloseChan:
				if ok {
					h.dataService.BatchUpsertDailyClose(ctx, objs)
				}
			case objs, ok := <-threePrimaryChan:
				if ok {
					h.dataService.BatchUpsertThreePrimary(ctx, objs)
				}
			case objs, ok := <-stakeConcentrationChan:
				if ok && len(*objs) > 0 {
					if val, ok := (*objs)[0].(*entity.StakeConcentration); ok {
						h.dataService.CreateStakeConcentration(ctx, &dto.CreateStakeConcentrationRequest{
							StockID:       val.StockID,
							Date:          val.Date,
							SumBuyShares:  val.SumBuyShares,
							SumSellShares: val.SumSellShares,
							AvgBuyPrice:   val.AvgBuyPrice,
							AvgSellPrice:  val.AvgSellPrice,
						})
						// refresh the concentration
						h.RefreshStakeConcentration(ctx, val.StockID, val.Date)
					}
				}
			}
		}
	}()
}

func (h *handlerImpl) generateJob(ctx context.Context, origin parser.Source, req *dto.DownloadRequest, respChan chan *[]interface{}) {
	for i := req.RewindLimit * -1; i <= 0; i++ {
		var date string
		switch origin {
		case parser.TwseDailyClose, parser.TwseThreePrimary:
			date = helper.GetDateFromOffset(i, helper.TwseDateFormat)
		case parser.TpexDailyClose, parser.TpexThreePrimary:
			date = helper.GetDateFromOffset(i, helper.TpexDateFormat)
		case parser.StakeConcentration:
			date = helper.GetDateFromOffset(i, helper.StakeConcentrationFormat)
		}

		if len(date) > 0 {
			if origin == parser.StakeConcentration {
				// align the date format to be 20220107, but remains the query date as 2022-01-07
				res, err := h.dataService.ListBackfillStakeConcentrationStockIDs(ctx, strings.ReplaceAll(date, "-", ""))
				if err != nil {
					log.Errorf("ListBackfillStakeConcentrationStockIDs error: %+v", err)
					continue
				}
				for _, id := range res {
					job := &downloadJob{
						ctx:       ctx,
						date:      date,
						stockId:   id,
						respChan:  respChan,
						rateLimit: req.RateLimit,
						origin:    origin,
					}

					select {
					case concurrent.JobQueue <- job:
					case <-ctx.Done():
						log.Debug("(BatchingDownload): generateJob goroutine exit!")
						return
					}
				}
			} else {
				job := &downloadJob{
					ctx:       ctx,
					date:      date,
					respChan:  respChan,
					rateLimit: req.RateLimit,
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
	}
	log.Debug("(BatchingDownload): all download jobs sent!")
}

func (job *downloadJob) Do() error {
	var c icrawler.ICrawler
	p := parser.New()
	var config parser.Config

	switch job.origin {
	case parser.TwseDailyClose:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 17,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		c.SetURL(icrawler.TwseDailyClose, job.date, icrawler.StockOnly)
	case parser.TpexDailyClose:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 17,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		c.SetURL(icrawler.TpexDailyClose, job.date)
	case parser.TwseThreePrimary:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 19,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		c.SetURL(icrawler.TwseThreePrimary, job.date, icrawler.StockOnly)
	case parser.TpexThreePrimary:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 24,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		c.SetURL(icrawler.TpexThreePrimary, job.date)
	case parser.TwseStockList:
		config = parser.Config{
			Capacity: 6,
			Type:     job.origin,
		}
		c = crawler.New(nil)
		c.SetURL(icrawler.TWSEStocks, "")
	case parser.TpexStockList:
		config = parser.Config{
			Capacity: 6,
			Type:     job.origin,
		}
		c = crawler.New(nil)
		c.SetURL(icrawler.TPEXStocks, "")
	case parser.StakeConcentration:
		config = parser.Config{
			ParseDay: &job.date,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.Concentration})
		c.SetURL(icrawler.StakeConcentration, job.date, job.stockId)
	default:
		return fmt.Errorf("no recognized job source being specified: %s", job.origin)
	}

	stream, err := c.Fetch(job.ctx)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
	}
	err = p.Parse(config, stream)
	if err != nil {
		sentry.CaptureException(err)
		return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
	}

	job.respChan <- p.Flush()

	// rate limit protection and context.cancel
	select {
	case <-time.After(time.Duration(job.rateLimit) * time.Millisecond):
	case <-job.ctx.Done():
		//log.Warn("(downloadJob) - context cancelled!")
	}

	return nil
}
