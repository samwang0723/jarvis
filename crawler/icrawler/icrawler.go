package icrawler

import (
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
	Fetch() (io.Reader, error)
	SetURL(template string, date string, options ...string)
}
