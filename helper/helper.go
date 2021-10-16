package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TB             = 1000000000000
	GB             = 1000000000
	MB             = 1000000
	KB             = 1000
	TimeZone       = "Asia/Taipei"
	TwseDateFormat = "20060102"
	TpexDateFormat = "2006/01/02"
)

func IsInteger(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func ToUint64(v string) uint64 {
	if i, err := strconv.ParseUint(v, 10, 64); err == nil {
		return i
	}

	return 0
}

func ToFloat32(v string) float32 {
	if f, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(f)
	}
	return 0
}

func ConvertDateStr(year int, month int, day int, format string) string {
	l, _ := time.LoadLocation(TimeZone)
	t := time.Now().In(l)
	t = t.AddDate(year, month, day)
	wkDay := t.Weekday()
	if wkDay == time.Saturday || wkDay == time.Sunday {
		return ""
	}
	s := t.Format(format)
	if format == TpexDateFormat {
		res := strings.Split(s, "/")
		year, _ := strconv.Atoi(res[0])
		s = fmt.Sprintf("%d/%s/%s", year-1911, res[1], res[2])
	}
	return s
}

func UnifiedDateStr(input string) string {
	if strings.Contains(input, "/") {
		i := strings.Split(input, "/")
		year, _ := strconv.Atoi(i[0])
		res := fmt.Sprintf("%d%s%s", year+1911, i[1], i[2])
		return res
	}
	return input
}

func DeserializeTime(date string, format string) (*time.Time, error) {
	d, err := time.Parse(format, date)
	if err != nil {
		return nil, err
	}
	l, _ := time.LoadLocation(TimeZone)
	d = d.In(l)
	return &d, nil
}

func ToDateStr(date time.Time, format string) string {
	return date.Format(format)
}

func ReadableSize(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	// Get whole number, and the remainder for decimals
	if length > TB {
		unit = "TB"
		i = length / TB
		remainder = length - (i * TB)
	} else if length > GB {
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
	} else if length > MB {
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
	} else if length > KB {
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
	} else {
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else if remainder > KB {
		width = 6
	} else {
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}
