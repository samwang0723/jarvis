package crawler

import (
	"io"
)

type Crawler interface {
	Fetch() (io.Reader, error)
	SetURL(template string, date string, queryType string)
}
