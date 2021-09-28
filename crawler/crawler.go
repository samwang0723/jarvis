package crawler

import (
	"io"
	"time"
)

type Crawler interface {
	Fetch() (io.Reader, error)
	SetURL(template string, date string, queryType string)
	SetRateLimit(sec time.Duration)
}
