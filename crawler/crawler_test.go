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

package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samwang0723/jarvis/crawler/icrawler"
	log "github.com/samwang0723/jarvis/logger"
	logtest "github.com/samwang0723/jarvis/logger/structured"
)

func setup() {
	logger := logtest.NullLogger()
	log.Initialize(logger)
}

func Test_SetURL(t *testing.T) {
	setup()
	tests := []struct {
		name     string
		template string
		date     string
		options  string
		want     string
	}{
		{
			name:     "Construct Twse dailyclose url",
			template: icrawler.TwseDailyClose,
			date:     "20211223",
			options:  icrawler.StockOnly,
			want:     "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=20211223&type=ALLBUT0999",
		},
		{
			name:     "Construct Tpex dailyclose url",
			template: icrawler.TpexDailyClose,
			date:     "20211223",
			want:     "http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=20211223&s=0,asc,0",
		},
		{
			name:     "Construct Twse threeprimary url",
			template: icrawler.TwseThreePrimary,
			date:     "20211223",
			options:  icrawler.StockOnly,
			want:     "http://www.tse.com.tw/fund/T86?response=csv&date=20211223&selectType=ALLBUT0999",
		},
		{
			name:     "Construct Twse stock list",
			template: icrawler.TWSEStocks,
			want:     "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2",
		},
		{
			name:     "Construct Tpex stock list",
			template: icrawler.TPEXStocks,
			want:     "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4",
		},
		{
			name:     "Construct stakeconcerntration",
			template: icrawler.StakeConcentration,
			date:     "20211223",
			options:  "1101",
			want:     "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco.djhtm?a=1101&e=20211223&f=20211223",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crawlerImpl{}
			if len(tt.options) > 0 {
				c.SetURL(tt.template, tt.date, tt.options)
			} else {
				c.SetURL(tt.template, tt.date)
			}

			if c.url != tt.want {
				t.Errorf("SetURL(%s, %s, %s) = %s, want %s", tt.template, tt.date, tt.options, c.url, tt.want)
			}
		})
	}
}

func Test_Fetch(t *testing.T) {
	setup()
	tests := []struct {
		name   string
		server *httptest.Server
		want   bool
	}{
		{
			name: "Regular http fetch",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Success"))

			})),
			want: false,
		},
		{
			name: "error fetching from server",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			})),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			c := &crawlerImpl{
				url:    tt.server.URL,
				client: tt.server.Client(),
			}
			_, err := c.Fetch(context.TODO())
			if (err != nil) != tt.want {
				t.Errorf("Fetch() = %v, want %v", err != nil, tt.want)
			}
		})
	}
}
