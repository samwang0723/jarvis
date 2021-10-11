package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"samwang0723/jarvis/crawler/icrawler"
	"samwang0723/jarvis/helper"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type crawlerImpl struct {
	url string
}

func New() icrawler.ICrawler {
	res := &crawlerImpl{}
	return res
}

func (c *crawlerImpl) SetURL(template string, date string, queryType string) {
	c.url = fmt.Sprintf(template, date, queryType)
}

func (c *crawlerImpl) Fetch() (io.Reader, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, fmt.Errorf("new fetch request initialize error: %v\n", err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
	}
	resp, err := client.Do(req)
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

	log.Printf("download completed (%s), URL: %s, Header: %v\n",
		helper.ReadableSize(len(csvfile), 2), c.url, resp.Header)
	raw := bytes.NewBuffer(csvfile)
	reader := transform.NewReader(raw, traditionalchinese.Big5.NewDecoder())

	return reader, nil
}
