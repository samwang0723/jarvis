package crawler

import (
	"io"
)

type ICrawler interface {
	Fetch() (io.Reader, error)
	SetURL(template string, date string, queryType string)
}
