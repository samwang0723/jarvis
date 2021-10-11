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
	TpexDailyClose   = "https://tpex.org.tw/web/stock/aftertrading/otc_quotes_no1430/stk_wn1430_result.php?l-zh-TW&o=csv&d=%s&se=EW&s=0,asc,0"
)

type ICrawler interface {
	Fetch() (io.Reader, error)
	SetURL(template string, date string, queryType string)
}
