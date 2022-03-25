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

	"github.com/getsentry/sentry-go"
	"github.com/samwang0723/jarvis/crawler"
	"github.com/samwang0723/jarvis/crawler/icrawler"
	"github.com/samwang0723/jarvis/crawler/proxy"
	"github.com/samwang0723/jarvis/parser"
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
		url := fmt.Sprintf(icrawler.TwseDailyClose, job.date)
		c.AppendURL(url)

	case parser.TpexDailyClose:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 17,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TpexDailyClose, job.date)
		c.AppendURL(url)

	case parser.TwseThreePrimary:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 19,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TwseThreePrimary, job.date)
		c.AppendURL(url)

	case parser.TpexThreePrimary:
		config = parser.Config{
			ParseDay: &job.date,
			Capacity: 24,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TpexThreePrimary, job.date)
		c.AppendURL(url)

	case parser.TwseStockList:
		config = parser.Config{
			Capacity: 6,
			Type:     job.origin,
		}
		c = crawler.New(nil)
		c.AppendURL(icrawler.TWSEStocks)

	case parser.TpexStockList:
		config = parser.Config{
			Capacity: 6,
			Type:     job.origin,
		}
		c = crawler.New(nil)
		c.AppendURL(icrawler.TPEXStocks)

	case parser.StakeConcentration:
		config = parser.Config{
			ParseDay: &job.date,
			Type:     job.origin,
		}
		c = crawler.New(&proxy.Proxy{Type: proxy.Concentration})

		// in order to get accurate data, we must query each page https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_6598_6.djhtm
		// as the top 15 brokers may different from day to day and not possible to store all detailed daily data
		indexes := []int{1, 2, 3, 4, 6}
		for _, idx := range indexes {
			c.AppendURL(fmt.Sprintf(icrawler.ConcentrationDays, job.stockId, idx))
		}

	default:
		return fmt.Errorf("no recognized job source being specified: %s", job.origin)
	}

	// looping to download all URLs
	for {
		urls := c.GetURLs()
		if len(urls) <= 0 {
			break
		}

		sourceURL, bytes, err := c.Fetch(job.ctx)
		if err != nil {
			sentry.CaptureException(err)
			return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
		}
		config.SourceURL = sourceURL
		err = p.Parse(config, bytes)
		if err != nil {
			sentry.CaptureException(err)
			return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
		}
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
