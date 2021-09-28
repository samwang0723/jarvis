package main

import (
	"log"
	"samwang0723/jarvis/crawler"
	"samwang0723/jarvis/helper"
	"samwang0723/jarvis/parser"
)

func main() {
	d := helper.GetDate(0, 0, -5)
	log.Println(d)

	var twse crawler.Crawler
	twse = &crawler.TwseStock{}
	twse.SetURL(crawler.DailyClose, d, crawler.StockOnly)
	io, err := twse.Fetch()
	if err != nil {
		log.Fatalf("Fetch Error: %s\n", err)
	}

	var handler parser.Parser
	handler = &parser.CsvHandler{Tag: d}
	data := map[string]interface{}{}
	handler.SetDataSource(data)
	config := parser.Config{
		StartInteger: true,
		Capacity:     17,
		Type:         parser.TwseDailyClose,
	}
	_, err = handler.Parse(config, io)
	if err != nil {
		log.Fatalf("Parse Error: %s\n", err)
	}

	for k, v := range data {
		log.Printf("%s: %v\n", k, v)
	}

}
