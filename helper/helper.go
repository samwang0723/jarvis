package helper

import (
	"fmt"
	"strconv"
	"time"
)

const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000
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

func GetDate(year int, month int, day int) string {
	t := time.Now()
	t = t.AddDate(year, month, day)
	return t.Format("20060102")
}

func ConvertDate(date string) (*time.Time, error) {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return nil, err
	}
	return &d, nil
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
