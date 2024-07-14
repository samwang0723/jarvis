package ericlagergren

import (
	"fmt"
)

type NoPlanError struct{}

func (e NoPlanError) Error() string {
	return "PlanScan did not find a plan"
}

type ScanError struct {
	Err error
}

func (e ScanError) Error() string {
	return fmt.Sprintf("Scan error: %s", e.Err)
}

type TypeAssertError struct{}

func (e TypeAssertError) Error() string {
	return "Type assertion failed"
}

type ScanNumericError struct {
	Val interface{}
}

func (e ScanNumericError) Error() string {
	return fmt.Sprintf("cannot scan %v into *decimal.Big", e.Val)
}

type ComposeError struct {
	Val interface{}
}

func (e ComposeError) Error() string {
	return fmt.Sprintf("fail to compose decimal.Big from %v", e.Val)
}
