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

	"github.com/samwang0723/jarvis/cache"
	"github.com/samwang0723/jarvis/concurrent"
	"github.com/samwang0723/jarvis/dto"
	"github.com/samwang0723/jarvis/entity"
	"github.com/samwang0723/jarvis/helper"
	log "github.com/samwang0723/jarvis/logger"
	"github.com/samwang0723/jarvis/parser"
)

const (
	StartCronjob = "START_CRON"
)

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
		// since we will have multiple daemonSet in nodes, need to make sure same cronjob
		// only running once at a time, here we use distrubted lock through Redis.
		lockAquired := cache.ObtainLock(cache.CronjobLock, 2*time.Minute)
		if lockAquired {
			h.BatchingDownload(longLiveCtx, &dto.DownloadRequest{
				RewindLimit: 0,
				RateLimit:   3000,
				Types:       req.Types,
			})
		} else {
			log.Error("CronDownload: Redis distributed lock obtain failed.")
		}
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

			concurrent.JobQueue <- job
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
					diff := []int32{0, 0, 0, 0, 0}
					stockId, date := "", ""
					for _, obj := range *objs {
						if val, ok := obj.(*entity.StakeConcentration); ok {
							switch val.HiddenField {
							case "1":
								h.dataService.CreateStakeConcentration(ctx, &dto.CreateStakeConcentrationRequest{
									StockID:       val.StockID,
									Date:          val.Date,
									SumBuyShares:  val.SumBuyShares,
									SumSellShares: val.SumSellShares,
									AvgBuyPrice:   val.AvgBuyPrice,
									AvgSellPrice:  val.AvgSellPrice,
								})
								stockId = val.StockID
								date = val.Date
								diff[0] = int32(val.SumBuyShares - val.SumSellShares)
							case "2":
								diff[1] = int32(val.SumBuyShares - val.SumSellShares)
							case "3":
								diff[2] = int32(val.SumBuyShares - val.SumSellShares)
							case "4":
								diff[3] = int32(val.SumBuyShares - val.SumSellShares)
							case "6":
								diff[4] = int32(val.SumBuyShares - val.SumSellShares)
							}
						}
					}
					// refresh the concentration
					h.RefreshStakeConcentration(ctx, &dto.RefreshStakeConcentrationRequest{StockID: stockId, Date: date, Diff: diff})
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
			var jobs []*downloadJob
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
					jobs = append(jobs, job)
				}
			} else {
				job := &downloadJob{
					ctx:       ctx,
					date:      date,
					respChan:  respChan,
					rateLimit: req.RateLimit,
					origin:    origin,
				}
				jobs = append(jobs, job)
			}

			for _, job := range jobs {
				concurrent.JobQueue <- job
			}
		}
	}
	log.Debug("(BatchingDownload): all download jobs sent!")
}
