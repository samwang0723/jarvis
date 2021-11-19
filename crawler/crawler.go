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
}

func New() icrawler.ICrawler {
	res := &crawlerImpl{
		url: "",
		client: &http.Client{
			Timeout: time.Second * 60,
		},
	}
	return res
}

func (c *crawlerImpl) SetURL(template string, date string, options ...string) {
	if len(options) > 0 {
		c.url = fmt.Sprintf(template, date, options[0])
	} else {
		c.url = fmt.Sprintf(template, date)
	}
}

func (c *crawlerImpl) Fetch(ctx context.Context) (io.Reader, error) {
	urlWithProxy := fmt.Sprintf("%s&url=%s", proxy.ProxyURI(), url.QueryEscape(c.url))
	req, err := http.NewRequest("GET", urlWithProxy, nil)
	if err != nil {
		return nil, fmt.Errorf("new fetch request initialize error: %v\n", err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
	}
	req = req.WithContext(ctx)
	log.Debugf("download started: %s\n", urlWithProxy)
	resp, err := (*c.client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch request error: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch status error: %v\n", resp.StatusCode)
	}

	csvfile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch unable to read body: %v\n", err)
	}

	log.Infof("download completed (%s), URL: %s", helper.ReadableSize(len(csvfile), 2), c.url)
	raw := bytes.NewBuffer(csvfile)
	reader := transform.NewReader(raw, traditionalchinese.Big5.NewDecoder())

	return reader, nil
}
