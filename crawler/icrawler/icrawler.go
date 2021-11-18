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
package icrawler

import (
	"context"
	"io"
)

const (
	All              = "ALL"
	StockOnly        = "ALLBUT0999"
	TwseDailyClose   = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=%s"
	TwseThreePrimary = "http://www.tse.com.tw/fund/T86?response=csv&date=%s&selectType=%s"
	OperatingDays    = "https://www.twse.com.tw/holidaySchedule/holidaySchedule?response=csv&queryYear=%d"
	TpexDailyClose   = "http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=%s&s=0,asc,0"
)

type ICrawler interface {
	Fetch(ctx context.Context) (io.Reader, error)
	SetURL(template string, date string, options ...string)
}
