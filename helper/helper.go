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
package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TB                       = 1000000000000
	GB                       = 1000000000
	MB                       = 1000000
	KB                       = 1000
	TimeZone                 = "Asia/Taipei"
	TwseDateFormat           = "20060102"
	TpexDateFormat           = "2006/01/02"
	StakeConcentrationFormat = "2006-01-02"
)

func IsInteger(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func ToInt64(v string) int64 {
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}

	return 0
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

func FormalizeValidTimeWithLocation(input time.Time, offset ...int) *time.Time {
	l, _ := time.LoadLocation(TimeZone)
	t := input.In(l)
	if len(offset) > 0 {
		t = t.AddDate(0, 0, offset[0])
	}

	// only within workday will be valid
	wkDay := t.Weekday()
	if wkDay == time.Saturday || wkDay == time.Sunday {
		return nil
	}
	return &t
}

func GetDateFromUTC(timestamp string, format string) string {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ""
	}
	t := FormalizeValidTimeWithLocation(time.Unix(i, 0))
	if t == nil {
		return ""
	}

	// Twse format: 20190213
	s := t.Format(format)

	switch format {
	case TpexDateFormat:
		// Tpex format: 108/02/06
		s = UnifiedDateFormatToTpex(s)
	}
	return s
}

func GetDateFromOffset(offset int, format string) string {
	t := FormalizeValidTimeWithLocation(time.Now(), offset)
	if t == nil {
		return ""
	}

	// Twse format: 20190213
	s := t.Format(format)

	switch format {
	case TpexDateFormat:
		// Tpex format: 108/02/06
		s = UnifiedDateFormatToTpex(s)
	}
	return s
}

func UnifiedDateFormatToTpex(input string) string {
	if strings.Contains(input, "/") {
		res := strings.Split(input, "/")
		year, _ := strconv.Atoi(res[0])
		return fmt.Sprintf("%d/%s/%s", year-1911, res[1], res[2])
	}
	return input
}

func UnifiedDateFormatToTwse(input string) string {
	if strings.Contains(input, "/") {
		i := strings.Split(input, "/")
		year, _ := strconv.Atoi(i[0])
		res := fmt.Sprintf("%d%s%s", year+1911, i[1], i[2])
		return res
	}
	return input
}

func GetReadableSize(length int, decimals int) (out string) {
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
