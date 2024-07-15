package helper_test

import (
	"flag"
	"math"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/google/go-cmp/cmp"
	"github.com/samwang0723/jarvis/internal/helper"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestBytes2String(t *testing.T) {
	t.Parallel()
	b := []byte("hello")
	expected := "hello"
	result := helper.Bytes2String(b)
	if result != expected {
		t.Errorf("Bytes2String(%v) = %v; want %v", b, result, expected)
	}
}

func TestRoundDecimal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    float32
		expected float32
	}{
		{1.5, 2},
		{2.3, 2},
		{2.7, 3},
	}

	for _, test := range tests {
		result := helper.RoundDecimal(test.input)
		if result != test.expected {
			t.Errorf("RoundDecimal(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestRoundDecimalTwo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    float32
		expected float32
	}{
		{1.234, 1.23},
		{1.235, 1.24},
		{1.236, 1.24},
		{2.345, 2.35},
		{2.344, 2.34},
		{0.0, 0.0},
		{-1.234, -1.23},
		{-1.235, -1.24},
	}

	for _, test := range tests {
		result := helper.RoundDecimalTwo(test.input)
		if result != test.expected {
			t.Errorf("RoundDecimalTwo(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestStringToFloat32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected float32
		hasError bool
	}{
		{"1.23", 1.23, false},
		{"abc", 0, true},
	}

	for _, test := range tests {
		result, err := helper.StringToFloat32(test.input)
		if (err != nil) != test.hasError {
			t.Errorf(
				"StringToFloat32(%v) error = %v; want error = %v",
				test.input,
				err,
				test.hasError,
			)
		}
		if result != test.expected {
			t.Errorf("StringToFloat32(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestUint64ToString(t *testing.T) {
	t.Parallel()
	num := uint64(1234567890)
	expected := "1234567890"
	result := helper.Uint64ToString(num)
	if result != expected {
		t.Errorf("Uint64ToString(%v) = %v; want %v", num, result, expected)
	}
}

func TestToday(t *testing.T) {
	t.Parallel()
	expected := time.Now().Format("20060102")
	result := helper.Today()
	if result != expected {
		t.Errorf("Today() = %v; want %v", result, expected)
	}
}

func TestRewindDate(t *testing.T) {
	t.Parallel()
	dateStr := "20231101"
	rewind := -1
	expected := "20231031"
	result := helper.RewindDate(dateStr, rewind)
	if result != expected {
		t.Errorf("RewindDate(%v, %v) = %v; want %v", dateStr, rewind, result, expected)
	}
}

func TestStringToUint64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected uint64
		hasError bool
	}{
		{"1234567890", 1234567890, false},
		{"0", 0, false},
		{"18446744073709551615", 18446744073709551615, false}, // Max uint64 value
		{"-1", 0, true},                                       // Invalid input
		{"abc", 0, true},                                      // Invalid input
	}

	for _, test := range tests {
		result, err := helper.StringToUint64(test.input)
		if (err != nil) != test.hasError {
			t.Errorf(
				"StringToUint64(%v) error = %v; want error = %v",
				test.input,
				err,
				test.hasError,
			)
		}
		if result != test.expected {
			t.Errorf("StringToUint64(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestUint64ToInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    uint64
		expected int
		hasError bool
	}{
		{1234567890, 1234567890, false},
		{0, 0, false},
		{uint64(math.MaxInt), math.MaxInt, false}, // Max int value
		{uint64(math.MaxInt) + 1, 0, true},        // Overflow case
	}

	for _, test := range tests {
		result, err := helper.Uint64ToInt(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Uint64ToInt(%v) error = %v; want error = %v", test.input, err, test.hasError)
		}
		if result != test.expected {
			t.Errorf("Uint64ToInt(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestFloat32ToDecimal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    float32
		expected string
	}{
		{123.456, "123.456"},
		{0.0, "0"},
		{-123.456, "-123.456"},
		{1.23456789, "1.2345679"}, // Precision of float32
	}

	for _, test := range tests {
		result := helper.Float32ToDecimal(test.input)
		if result.String() != test.expected {
			t.Errorf(
				"Float32ToDecimal(%v) = %v; want %v",
				test.input,
				result.String(),
				test.expected,
			)
		}
	}
}

func TestDecimalToFloat32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected float32
	}{
		{"123.456", 123.456},
		{"0", 0.0},
		{"-123.456", -123.456},
		{"1.23456789", 1.2345679}, // Precision of float32
	}

	for _, test := range tests {
		var d decimal.Big
		if _, ok := d.SetString(test.input); !ok {
			t.Errorf("Invalid decimal string: %v", test.input)
			continue
		}
		result := helper.DecimalToFloat32(d)
		if result != test.expected {
			t.Errorf("DecimalToFloat32(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    map[string]int
		expected []string
	}{
		{map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "b", "c"}},
		{map[string]int{}, []string{}},
		{map[string]int{"one": 1}, []string{"one"}},
	}

	for _, test := range tests {
		result := helper.Keys(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Keys(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestSliceToMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    []int
		fn       func(int) string
		expected map[string]int
	}{
		{
			[]int{1, 2, 3},
			func(i int) string { return strconv.Itoa(i) },
			map[string]int{"1": 1, "2": 2, "3": 3},
		},
		{
			[]int{},
			func(i int) string { return strconv.Itoa(i) },
			map[string]int{},
		},
		{
			[]int{10, 20, 30},
			func(i int) string { return "key" + strconv.Itoa(i) },
			map[string]int{"key10": 10, "key20": 20, "key30": 30},
		},
	}

	for _, test := range tests {
		result := helper.SliceToMap(test.input, test.fn)
		if !cmp.Equal(result, test.expected) {
			t.Errorf("Query diff = %v", cmp.Diff(test.input, test.expected))
		}
	}
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		str      string
		list     []string
		expected bool
	}{
		{"a", []string{"a", "b", "c"}, true},
		{"d", []string{"a", "b", "c"}, false},
		{"", []string{"a", "b", "c"}, false},
		{"a", []string{}, false},
	}

	for _, test := range tests {
		result := helper.StringInSlice(test.str, test.list)
		if result != test.expected {
			t.Errorf(
				"StringInSlice(%v, %v) = %v; want %v",
				test.str,
				test.list,
				result,
				test.expected,
			)
		}
	}
}
