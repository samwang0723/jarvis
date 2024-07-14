package ericlagergren_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/samwang0723/jarvis/internal/db/pginit/ext/decimal/ericlagergren"
)

func TestNoPlanError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ericlagergren.NoPlanError{},
			expectedString: "PlanScan did not find a plan",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestScanError(t *testing.T) {
	t.Parallel()

	err := errors.New("test")

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ericlagergren.ScanError{Err: err},
			expectedString: fmt.Sprintf("Scan error: %s", err),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestTypeAssertError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ericlagergren.TypeAssertError{},
			expectedString: "Type assertion failed",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestScanNumericError(t *testing.T) {
	t.Parallel()

	val := "NULL"

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ericlagergren.ScanNumericError{Val: val},
			expectedString: fmt.Sprintf("cannot scan %v into *decimal.Big", val),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestComposeError(t *testing.T) {
	t.Parallel()

	val := "NULL"

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ericlagergren.ComposeError{Val: val},
			expectedString: fmt.Sprintf("fail to compose decimal.Big from %v", val),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}
