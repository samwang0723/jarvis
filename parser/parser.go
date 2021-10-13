package parser

import "io"

const (
	TwseDailyClose = iota
	TwseThreePrimary
	TpexDailyClose
)

type IParser interface {
	Parse(config Config, in io.Reader) error
	Flush() *[]interface{}
}

type parserImpl struct {
	result *[]interface{}
}

type Config struct {
	ParseDay *string
	Capacity int
	Type     int
}

func New() IParser {
	res := &parserImpl{
		result: &[]interface{}{},
	}
	return res
}
