package parser

import "io"

const (
	TwseDailyClose = iota
	TwseThreePrimary
	TpexDailyClose
)

type IParser interface {
	Parse(config Config, in io.Reader) (*[]interface{}, error)
	SetDataSource(source *[]interface{})
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
	res := &parserImpl{}
	return res
}
