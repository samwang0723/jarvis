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

package parser

import (
	"strings"
	"testing"
)

func Test_parseCsv(t *testing.T) {
	wrongCsv := `
"110年11月30日每日收盤行情(全部(不含權證、牛熊證))"
"(元,股)",,,,,,,,,,,"(元,交易單位)",,,,,
"證券代號","證券名稱","成交股數","成交筆數","成交金額","開盤價","最高價","最低價","收盤價","漲跌(+/-)","漲跌價差","最後揭示買價","最後揭示買量","最後揭示賣價","最後揭示賣量","本益比",
="020008","元大特股高息N","19,000","6","240,690","12.67","12.68","12.66","12.68","+","0.02","12.65","100","12.66","100","0.00",
="020011","統一微波高息20N","5,000","3","33,990","6.80","6.80","6.79","6.79","+","0.01","6.76","493","6.77","281","0.00",
="020012","富邦行動通訊N","33,000","6","214,190","6.50","6.50","6.48","6.48","+","0.07","6.40","100","6.41","100","0.00",
="020015","統一MSCI美低波N","0","0","0","--","--","--","--"," ","0.00","15.37","499","15.39","280","0.00",
="020016","統一MSCI美科技N","5,000","3","109,650","21.93","21.93","21.93","21.93","+","0.31","21.84","492","21.86","281","0.00",
="020018","統一價值成長30N","28,000","6","517,100","18.50","18.50","18.37","18.37","+","0.19","18.34","492","18.36","286","0.00",
="020019","統一特選台灣5GN","12,000","7","179,580","14.89","15.01","14.88","14.88","+","0.10","14.86","499","14.88","283","0.00",
="02001L","富邦蘋果正二N","31,000","2","317,420","10.22","10.24","10.22","10.24","+","0.17","10.13","10","10.14","10","0.00",
="02001R","富邦蘋果反一N","2,000","1","6,160","3.08","3.08","3.08","3.08","-","0.02","3.09","10","3.10","10","0.00",
="020020","元大台股領航N","314,000","97","2,298,580","7.28","7.38","7.27","7.28","+","0.01","7.27","100","7.28","90","0.00",
	`
	correctCsv := `
"110年11月30日 價格指數(臺灣證券交易所)"
"指數","收盤指數","漲跌(+/-)","漲跌點數","漲跌百分比(%)","特殊處理註記",
"寶島股價指數","19,959.18","+","123.68","0.62","",
"發行量加權股價指數","17,427.76","+","99.67","0.58","",
"臺灣公司治理100指數","10,114.29","+","48.00","0.48","",
"臺灣50指數","13,590.78","+","42.48","0.31","",
"臺灣50權重上限30%指數","12,788.54","+","30.99","0.24","",
"臺灣中型100指數","14,511.81","+","124.26","0.86","",
"臺灣資訊科技指數","23,903.42","+","127.50","0.54","",
"臺灣發達指數","10,920.75","+","11.37","0.10","",
"臺灣高股息指數","7,569.09","-","3.64","-0.05","",
"110年11月30日 大盤統計資訊"
"成交統計","成交金額(元)","成交股數(股)","成交筆數",
"1.一般股票","438,925,293,798","6,236,206,018","2,298,510",
"2.台灣存託憑證","671,179,843","45,875,365","14,247",
"3.受益憑證","0","0","0",
"4.ETF","8,515,694,192","469,270,804","103,488",
"5.受益證券","73,294,280","7,231,809","90",
"6.變更交易股票","15,335,682","3,364,840","1,079",
"7.認購(售)權證","2,651,260,350","2,284,241,000","94,828",
"8.轉換公司債","0","0","0",
"9.附認股權特別股","0","0","0",
"10.附認股權公司債","0","0","0",
"11.認股權憑證","0","0","0",
"12.公司債","0","0","0",
"13.ETN","59,734,010","11,382,000","1,360",
"14.創新板股票","0","0","0",
"15.創新板-變更交易方法股票","0","0","0",
"證券合計(1+6+14+15)","438,940,629,480","6,239,570,858","2,299,589",
"總計(1~15)","450,911,792,155","9,057,571,836","2,513,602",
"漲跌證券數合計"
"類型","整體市場","股票",
"上漲(漲停)","7,206(20)","592(13)",
"下跌(跌停)","3,060(59)","276(1)",
"持平","609","74",
"未成交","11,239","4",
"無比價","2,472","7",
"備註:"
""漲跌價差"為當日收盤價與前一日收盤價比較。",
""無比價"含前一日無收盤價、當日除權、除息、新上市、恢復交易者。",
"外幣成交值係以本公司當日下午3時30分公告匯率換算後加入成交金額。<br>公告匯率請參考本公司首頁>產品與服務>交易系統>雙幣ETF專區>代號對應及每日公告匯率。",

"110年11月30日每日收盤行情(全部(不含權證、牛熊證))"
"(元,股)",,,,,,,,,,,"(元,交易單位)",,,,,
"證券代號","證券名稱","成交股數","成交筆數","成交金額","開盤價","最高價","最低價","收盤價","漲跌(+/-)","漲跌價差","最後揭示買價","最後揭示買量","最後揭示賣價","最後揭示賣量","本益比",
="0050","元大台灣50","10,539,304","8,776","1,460,465,604","138.80","139.60","137.85","138.00","-","0.15","138.00","98","138.10","20","0.00",
="0051","元大中型100","148,884","108","8,872,890","59.70","59.70","59.45","59.45","+","0.50","59.40","2","59.55","57","0.00",
="0052","富邦科技","538,230","199","69,249,413","128.85","129.25","127.50","127.70","+","0.20","127.70","22","127.95","3","0.00",
="0053","元大電子","11,026","12","742,957","67.70","67.70","67.20","67.20","+","0.30","66.95","2","67.45","11","0.00",
"020004","兆豐電菁英30N","0","0","0","--","--","--","--"," ","0.00","32.69","499","32.70","499","0.00",
="020006","永昌中小300N","0","0","0","--","--","--","--"," ","0.00","34.91","100","34.93","100","0.00",
="020007","凱基臺灣500N","0","0","0","--","--","--","--"," ","0.00","34.45","218","34.51","218","0.00",
="020008","元大特股高息N","19,000","6","240,690","12.67","12.68","12.66","12.68","+","0.02","12.65","100","12.66","100","0.00",
="020011","統一微波高息20N","5,000","3","33,990","6.80","6.80","6.79","6.79","+","0.01","6.76","493","6.77","281","0.00",
="020012","富邦行動通訊N","33,000","6","214,190","6.50","6.50","6.48","6.48","+","0.07","6.40","100","6.41","100","0.00",
="020015","統一MSCI美低波N","0","0","0","--","--","--","--"," ","0.00","15.37","499","15.39","280","0.00",
="020016","統一MSCI美科技N","5,000","3","109,650","21.93","21.93","21.93","21.93","+","0.31","21.84","492","21.86","281","0.00",
="020018","統一價值成長30N","28,000","6","517,100","18.50","18.50","18.37","18.37","+","0.19","18.34","492","18.36","286","0.00",
="020019","統一特選台灣5GN","12,000","7","179,580","14.89","15.01","14.88","14.88","+","0.10","14.86","499","14.88","283","0.00",
="02001L","富邦蘋果正二N","31,000","2","317,420","10.22","10.24","10.22","10.24","+","0.17","10.13","10","10.14","10","0.00",
="02001R","富邦蘋果反一N","2,000","1","6,160","3.08","3.08","3.08","3.08","-","0.02","3.09","10","3.10","10","0.00",
="020020","元大台股領航N","314,000","97","2,298,580","7.28","7.38","7.27","7.28","+","0.01","7.27","100","7.28","90","0.00",
="020022","元大電動車N","649,000","136","3,529,900","5.45","5.47","5.41","5.41","+","0.02","5.41","99","5.42","104","0.00",
="020028","元大特選電動車N","947,000","174","4,966,000","5.25","5.27","5.23","5.23","+","0.02","5.22","522","5.23","493","0.00",
="020029","元大ESG高股息N","887,000","192","4,654,930","5.25","5.29","5.21","5.21"," ","0.00","5.20","623","5.21","503","0.00",
="020030","統一智慧電動車N","8,423,000","722","42,478,100","5.04","5.07","5.00","5.02","+","0.07","5.01","517","5.02","437","0.00",
"1101","台泥","34,666,716","10,910","1,611,298,375","46.50","47.30","46.00","46.00","-","0.50","46.00","489","46.10","1","13.49",
"1101B","台泥乙特","1,606","8","83,651","52.30","52.30","52.30","52.30"," ","0.00","51.80","7","52.30","31","0.00",
"1102","亞泥","21,163,999","6,038","896,536,212","42.60","43.15","42.00","42.00","-","0.45","42.00","524","42.15","40","9.63",
"1103","嘉泥","747,424","376","15,382,013","20.65","20.75","20.30","20.75","+","0.45","20.55","1","20.75","16","7.31",
"1104","環泥","351,210","165","7,372,588","21.00","21.10","20.85","21.00","+","0.10","21.00","10","21.05","11","12.43",
"1108","幸福","243,041","148","2,783,320","11.45","11.55","11.40","11.40","-","0.05","11.40","8","11.50","1","19.66",
"1109","信大","398,759","217","8,019,269","20.10","20.25","20.00","20.00","-","0.05","20.00","35","20.05","3","7.07",
"1110","東泥","777,054","374","15,452,420","20.20","20.25","19.65","19.65","-","0.20","19.60","24","19.65","7","81.88",
"1201","味全","1,031,795","429","22,941,554","22.35","22.35","22.05","22.30","+","0.15","22.20","14","22.30","19","18.90",
"1203","味王","5,068","21","167,213","33.00","33.00","33.00","33.00","-","0.20","33.00","16","33.20","4","18.03",
"1210","大成","1,613,323","1,661","84,212,936","52.30","52.40","52.00","52.00","-","0.20","52.00","369","52.10","10","18.12",
"1213","大飲","1,123","3","9,389","8.41","8.41","8.41","8.41","+","0.01","8.26","3","8.49","2","0.00",
"1215","卜蜂","625,137","613","50,395,826","80.00","81.00","79.70","80.90","+","0.90","80.20","5","80.90","2","14.63",
"1216","統一","22,691,588","5,615","1,498,351,003","67.10","67.50","65.40","65.40","-","1.90","65.40","1,344","65.60","81","18.53",
"1217","愛之味","1,056,894","351","10,879,197","10.30","10.35","10.25","10.30","+","0.05","10.25","200","10.35","107","20.20",
"1218","泰山","215,443","187","5,766,427","26.75","26.85","26.65","26.75"," ","0.00","26.75","14","26.85","10","21.75",
"1219","福壽","219,202","143","4,496,872","20.50","20.65","20.45","20.50"," ","0.00","20.45","12","20.55","5","16.67",
"1220","台榮","130,479","78","2,003,202","15.50","15.55","15.30","15.30","-","0.05","15.30","47","15.40","10","13.19",
"1225","福懋油","91,254","60","4,578,133","50.40","50.50","49.95","50.30","-","0.10","50.30","1","50.40","1","22.97",
"1227","佳格","803,769","723","41,180,225","51.30","51.50","51.00","51.00","-","0.30","51.00","158","51.40","1","17.53",
"1229","聯華","1,446,377","1,427","82,233,893","56.20","57.50","56.20","57.50","+","1.40","57.10","3","57.50","61","16.34",
"1231","聯華食","94,050","174","6,411,997","67.80","68.60","67.60","68.00","+","0.40","68.00","8","68.10","3","17.26",
"1232","大統益","75,055","176","12,083,062","161.00","162.00","160.00","160.50","+","0.50","160.50","2","161.50","1","19.48",
"1233","天仁","7,102","10","238,880","33.60","33.65","33.60","33.65","-","0.50","33.75","1","33.90","1","73.15",
"1234","黑松","72,309","60","2,524,998","34.75","35.00","34.75","34.80","+","0.05","34.80","9","34.85","4","19.77",
"1235","興泰","14,203","24","631,368","44.25","44.50","44.25","44.50"," ","0.00","44.50","2","44.75","1","61.81",
"1236","宏亞","901,794","680","19,374,087","20.90","21.85","20.10","21.75","+","1.85","21.70","7","21.75","20","0.00",
"1256","鮮活果汁-KY","25,199","563","9,023,684","359.00","359.00","356.50","356.50","-","1.50","357.00","1","358.00","3","14.03",
"1301","台塑","17,880,875","4,386","1,839,713,922","102.50","104.00","102.00","102.50","-","0.50","102.50","1,024","103.00","366","10.01",
`

	tests := []struct {
		name    string
		content string
		want    int
	}{
		{
			name:    "normal csv",
			content: correctCsv,
			want:    29,
		},
		{
			name:    "wrong csv",
			content: wrongCsv,
			want:    0,
		},
	}

	for _, tt := range tests {
		tt := tt
		date := "20211130"
		conf := Config{
			ParseDay: &date,
			Capacity: 17,
			Type:     TwseDailyClose,
		}

		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			in := strings.NewReader(tt.content)
			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.parseCsv(conf, in)

			if got := len(*res.result); got != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", got, tt.want)
			}
		})
	}

}