package sqlc

import (
	"errors"
	"fmt"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

type NotEnoughBalanceError struct {
	Err error
}

func (e *NotEnoughBalanceError) Error() string {
	return fmt.Sprintf("not enough balance", e.Err)
}

func (e *NotEnoughBalanceError) Unwrap() error {
	return e.Err
}

type RecordNotFoundError struct {
	Err error
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("record not found: %s", e.Err)
}

func (e *RecordNotFoundError) Unwrap() error {
	return e.Err
}

func newRecordNotFoundError(err error) error {
	return &RecordNotFoundError{Err: err}
}

func IsRecordNotFoundError(err error) bool {
	var e *RecordNotFoundError

	return errors.As(err, &e)
}

type TypeMismatchError struct {
	expect any
	got    any
}

func (tme *TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch, expect %t, got %t", tme.expect, tme.got)
}

type DuplicatedRecordError string

func (e DuplicatedRecordError) Error() string {
	return fmt.Sprintf("duplicated record error: %s", string(e))
}

type UnsupportedTradeTypeError string

func (e UnsupportedTradeTypeError) Error() string {
	return fmt.Sprintf("trade type unsupported: %s", string(e))
}

type DuplicatedEventError struct {
	event eventsourcing.Event
}

func (e *DuplicatedEventError) Error() string {
	return fmt.Sprintf("duplicated event: %s", string(e.event.EventType()))
}

type UnexpectedNadexMemberStatusError string

func (e UnexpectedNadexMemberStatusError) Error() string {
	return fmt.Sprintf("unexpected nadex member status: %s", string(e))
}
