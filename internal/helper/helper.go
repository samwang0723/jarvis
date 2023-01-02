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
	"math"
	"os"
	"reflect"
	"unsafe"
)

const (
	TimeZone = "Asia/Taipei"
)

func GetCurrentEnv() string {
	env := os.Getenv("ENVIRONMENT")
	output := "dev"
	switch env {
	case "development":
		output = "dev"
	case "staging":
		output = "staging"
	case "production":
		output = "prod"
	}
	return output
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
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

func RoundUpDecimalTwo(x float32) float32 {
	return float32(math.Ceil(float64(x)*100) / 100)
}
