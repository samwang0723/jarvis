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
	"math"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

const (
	TimeZone    = "Asia/Taipei"
	baseDecimal = 100
	floatFormat = 32
	uintFormat  = 64
	uintBase    = 10
)

func GetCurrentEnv() string {
	env := os.Getenv("ENVIRONMENT")
	output := "dev"

	switch env {
	case "local":
		output = "local"
	case "development":
		output = "dev"
	case "staging":
		output = "staging"
	case "production":
		output = "prod"
	}

	return output
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Convert a slice or array of a specific type to array of interface{}
func CastInterfaceSlice(s interface{}) *[]interface{} {
	v := reflect.ValueOf(s)
	// There is no need to check, we want to panic if it's not slice or array
	intf := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		intf[i] = v.Index(i).Interface()
	}

	return &intf
}

func RoundDecimal(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func RoundUpDecimalTwo(x float32) float32 {
	return float32(math.Ceil(float64(x)*baseDecimal) / baseDecimal)
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
