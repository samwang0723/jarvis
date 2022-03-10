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
	"time"

	"github.com/samwang0723/jarvis/concurrent"
	"github.com/samwang0723/jarvis/dto"
	"github.com/samwang0723/jarvis/entity"
	"github.com/samwang0723/jarvis/helper"
	log "github.com/samwang0723/jarvis/logger"
)

type refreshJob struct {
	ctx       context.Context
	date      string
	stockId   string
	respChan  chan *entity.StakeConcentration
	handler   *handlerImpl
	rateLimit int
}

func (h *handlerImpl) RefreshConcentration(ctx context.Context, rewindLimit int32) error {
	respChan := make(chan *entity.StakeConcentration)
	stocks, err := h.ListStock(ctx, &dto.ListStockRequest{
		Offset: 0,
		Limit:  2000,
		SearchParams: &dto.ListStockSearchParams{
			Country: "TW",
		},
	})
	log.Info("starting stake concentration refresh!")
	if err != nil {
		return fmt.Errorf("failed to retrieve stock list: %s", err)
	}

	go h.generateRefreshJob(ctx, stocks.Entries, rewindLimit, respChan)

	go func() {
		buffer := &[]interface{}{}
		flushSize := 50

		// make sure to flush remaining buffer
		defer func() {
			if len(*buffer) > 0 {
				h.dataService.BatchUpdateStakeConcentration(ctx, buffer)
			}
		}()

		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(2 * time.Hour):
				log.Warn("(RefreshConcentration): timeout")
				return
			case <-ctx.Done():
				log.Warn("(RefreshConcentration): context cancel")
				return
			case obj, ok := <-respChan:
				if ok && obj != nil {
					if len(*buffer) >= flushSize {
						// refresh the concentration
						h.dataService.BatchUpdateStakeConcentration(ctx, buffer)
						buffer = &[]interface{}{}
						log.Info("batch update stake concentration!")
					} else {
						*buffer = append(*buffer, obj)
					}
				}
			}
		}
	}()

	return nil
}

func (h *handlerImpl) generateRefreshJob(ctx context.Context, stocks []*entity.Stock, rewindLimit int32, respChan chan *entity.StakeConcentration) {
	for i := rewindLimit * -1; i <= 0; i++ {
		date := helper.GetDateFromOffset(i, helper.TwseDateFormat)
		if len(date) <= 0 || !h.dataService.HasStakeConcentration(ctx, date) {
			continue
		}
		for _, stock := range stocks {
			job := &refreshJob{
				ctx:       ctx,
				date:      date,
				stockId:   stock.StockID,
				respChan:  respChan,
				handler:   h,
				rateLimit: 10000,
			}
			select {
			case concurrent.JobQueue <- job:
			case <-ctx.Done():
				log.Debug("(RefreshConcentration): generateJob goroutine exit!")
				return
			}
		}
	}
}

func (job *refreshJob) Do() error {
	c := job.handler.calculateConcentration(job.ctx, job.stockId, job.date)
	if c == nil {
		return fmt.Errorf("failed to calculate concentration of %s, %s", job.stockId, job.date)
	}

	job.respChan <- c

	// rate limit protection and context.cancel
	select {
	case <-time.After(time.Duration(job.rateLimit) * time.Millisecond):
	case <-job.ctx.Done():
		//log.Warn("(refreshJob) - context cancelled!")
	}
	return nil
}
