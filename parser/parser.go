package parser

import "io"

const (
	TwseDailyClose = iota
	TwseThreePrimary
)

type Parser interface {
	Parse(config Config, in io.Reader) (map[string]interface{}, error)
	SetDataSource(map[string]interface{})
}

type Config struct {
	StartInteger bool
	Capacity     int
	Type         int
}
