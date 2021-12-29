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
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "samwang0723/jarvis/logger"

	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/crawler/proxy"
	"samwang0723/jarvis/helper"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type crawlerImpl struct {
	url    string
	client *http.Client
	proxy  *proxy.Proxy
}

func New(p *proxy.Proxy) icrawler.ICrawler {
	res := &crawlerImpl{
		url: "",
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		proxy: p,
	}
	return res
}

func (c *crawlerImpl) SetURL(template string, date string, options ...string) {
	if len(options) > 0 {
		if template == icrawler.StakeConcentration {
			c.url = fmt.Sprintf(template, options[0], date, date)
		} else {
			c.url = fmt.Sprintf(template, date, options[0])
		}
	} else {
		if len(date) == 0 {
			c.url = template
		} else {
			c.url = fmt.Sprintf(template, date)
		}
	}
}

func (c *crawlerImpl) Fetch(ctx context.Context) (io.Reader, error) {
	uri := c.url
	if c.proxy != nil {
		uri = fmt.Sprintf("%s&url=%s", c.proxy.URI(), url.QueryEscape(c.url))
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("new fetch request initialize error: %v", err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
		// It is important to close the connection otherwise fd count will overhead
		"Connection": []string{"close"},
	}
	req = req.WithContext(ctx)
	log.Debugf("download started: %s", uri)
	resp, err := (*c.client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch request error: %v, url: %s", err, c.url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch status error: %v, url: %s", resp.StatusCode, c.url)
	}

	// copy stream from response body, although it consumes memory but
	// better helps on concurrent handling in goroutine.
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch unable to read body: %v, url: %s", err, c.url)
	}

	log.Debugf("download completed (%s), URL: %s", helper.GetReadableSize(len(f), 2), c.url)
	raw := bytes.NewBuffer(f)
	reader := transform.NewReader(raw, traditionalchinese.Big5.NewDecoder())

	return reader, nil
}
