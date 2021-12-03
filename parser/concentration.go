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
	"fmt"
	"io"
	"samwang0723/jarvis/entity"
	"samwang0723/jarvis/helper"
	"strings"

	"golang.org/x/net/html"
)

const (
	sumBuy  = "合計買超張數"
	sumSell = "合計賣超張數"
	avgBuy  = "平均買超成本"
	avgSell = "平均賣超成本"
)

func (p *parserImpl) parseConcentration(config Config, in io.Reader) error {
	doc, err := html.Parse(in)
	if err != nil {
		return fmt.Errorf("failed to parse concentration: %s", err)
	}

	var concentration *entity.StakeConcentration

	// parse the header of stockID
	if title, ok := getHtmlTitle(doc); ok {
		t := strings.Split(title, "-")
		concentration = &entity.StakeConcentration{
			Date:    *config.ParseDay,
			StockID: strings.TrimSpace(t[1]),
		}
	}
	if concentration == nil {
		return fmt.Errorf("failed to parser concentration StockID: %+v", nil)
	}

	// parser the content of stake concentration
	var tag string
	for i := 1; i <= 2; i++ {
		row := getElementById(doc, "oScrollFoot", i)
		if row != nil && row.Data == "tr" {
			for c := row.FirstChild; c != nil; c = c.NextSibling {
				for d := c.FirstChild; d != nil && d.Type == html.TextNode; d = d.NextSibling {
					t := strings.Replace(d.Data, ",", "", -1)
					if helper.ToUint64(t) > 0 {
						switch tag {
						case sumBuy:
							concentration.SumBuyShares = helper.ToUint64(t)
						case sumSell:
							concentration.SumSellShares = helper.ToUint64(t)
						}
					} else if helper.ToFloat32(t) > 0 {
						switch tag {
						case avgBuy:
							concentration.AvgBuyPrice = helper.ToFloat32(t)
						case avgSell:
							concentration.AvgSellPrice = helper.ToFloat32(t)
						}
					} else {
						tag = strings.TrimSpace(d.Data)
					}
				}
			}
		}
	}

	*p.result = append(*p.result, concentration)

	return nil
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, id string, target int, cursor *int) *html.Node {
	if checkId(n, id) {
		*cursor++
		if *cursor == target {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, id, target, cursor)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementById(n *html.Node, id string, target int) *html.Node {
	cursor := 0
	return traverse(n, id, target, &cursor)
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverseTitle(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverseTitle(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func getHtmlTitle(n *html.Node) (string, bool) {
	return traverseTitle(n)
}