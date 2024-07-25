package domain

import (
	"flag"
	"os"
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expect  *Realtime
		name    string
		jsonStr string
	}{
		{
			name: "realtime json unmarshal successfully",
			jsonStr: `{"msgArray":[{"tv":"1","ps":"-","pz":"-","bp":"0","a":"484.0000_484.5000_485.0000_485.5000_486.0000_",
			"b":"483.5000_483.0000_482.5000_482.0000_481.5000_","c":"2330","d":"20230111","ch":"2330.tw",
			"tlong":"1673400815000","f":"68_53_189_135_407_","ip":"0","g":"467_617_66_126_97_","mt":"793167",
			"h":"488.0000","i":"24","it":"12","l":"482.0000","n":"台積電","o":"487.0000","p":"0","ex":"tse",
			"s":"1","t":"09:33:35","u":"534.0000","v":"5761","w":"437.5000","nf":"台灣積體電路製造股份有限公司",
			"y":"486.0000","z":"484.0000","ts":"0"}],"referer":"","userDelay":5000,"rtcode":"0000",
			"queryTime":{"sysDate":"20230111","stockInfoItem":2179,"stockInfo":964720,"sessionStr":"UserSession","sysTime":"09:33:38",
			"showChart":false,"sessionFromTime":1673400806975,"sessionLatestTime":1673400806975},"rtmessage":"OK",
			"exKey":"if_tse_2330.tw_zh-tw.null","cachedAlive":4251}`,
			expect: &Realtime{
				StockID:   "2330",
				Date:      "20230111",
				Open:      float32(487.0),
				Close:     float32(484.0),
				High:      float32(488.0),
				Low:       float32(482.0),
				Volume:    uint64(5761),
				ParseTime: "09:33:35",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			realtime := &Realtime{}
			realtime.UnmarshalJSON([]byte(tt.jsonStr))

			if *tt.expect != *realtime {
				t.Errorf("expect %+v, got %+v", tt.expect, realtime)
			}
		})
	}
}
