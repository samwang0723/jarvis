// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package helper

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"unsafe"

	"github.com/ericlagergren/decimal"
)

const (
	TimeZone    = "Asia/Taipei"
	baseDecimal = 100
	floatFormat = 32
	uintFormat  = 64
	uintBase    = 10
)

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func RoundDecimal(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func RoundDecimalTwo(x float32) float32 {
	return float32(math.Round(float64(x)*baseDecimal) / baseDecimal)
}

func StringToFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, floatFormat)
	if err != nil {
		return 0, err
	}

	return float32(f), nil
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func StringToUint64(s string) (uint64, error) {
	f, err := strconv.ParseUint(s, uintBase, uintFormat)
	if err != nil {
		return 0, err
	}

	return f, nil
}

func Uint64ToInt(u uint64) (int, error) {
	// Check if the value is within the range of int
	if u > math.MaxInt {
		return 0, fmt.Errorf("overflow: cannot convert %d to int", u)
	}
	return int(u), nil
}

func Float32ToDecimal(f float32) decimal.Big {
	// Convert float32 to string
	fStr := strconv.FormatFloat(float64(f), 'f', -1, 32)
	// Create a new decimal.Big
	var d decimal.Big
	// Set the value of decimal.Big from the string
	if _, ok := d.SetString(fStr); !ok {
		return decimal.Big{}
	}
	return d
}

func DecimalToFloat32(d decimal.Big) float32 {
	// Convert decimal.Big to string
	dStr := d.String()
	// Convert string to float32
	f, err := strconv.ParseFloat(dStr, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

func Today() string {
	return time.Now().Format("20060102")
}

func RewindDate(dateStr string, rewind int) string {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return ""
	}

	return date.AddDate(0, 0, rewind).Format("20060102")
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func SliceToMap[K comparable, V any](source []V, fn func(in V) K) map[K]V {
	output := map[K]V{}
	for _, s := range source {
		output[fn(s)] = s
	}

	return output
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}
