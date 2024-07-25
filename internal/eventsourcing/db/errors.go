package db

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

type AggregateNotFoundError struct {
	err         error
	aggregateID uuid.UUID
}

func (anfe AggregateNotFoundError) Error() string {
	return fmt.Sprintf("aggregate not found: %s (id %s)", anfe.err, anfe.aggregateID)
}

type LoadEventError struct {
	err error
}

func (lee LoadEventError) Error() string {
	return fmt.Sprintf("failed to load event: %s", lee.err)
}

type ApplyEventError struct {
	err         error
	aggregateID uuid.UUID
}

func (aee ApplyEventError) Error() string {
	return fmt.Sprintf("apply event error: %s (aggregate_id %s)", aee.err, aee.aggregateID)
}

type EventPublishError struct {
	err   error
	event eventsourcing.Event
}

func (epe EventPublishError) Error() string {
	return fmt.Sprintf("failed to publish event : %s (aggregate_id %s, event %s)",
		epe.err, epe.event.GetAggregateID(), epe.event.EventType())
}

type EventProjectorError struct {
	err   error
	event eventsourcing.Event
}

func (epe EventProjectorError) Error() string {
	return fmt.Sprintf("failed to project event : %s (aggregate_id %s, event %s)",
		epe.err, epe.event.GetAggregateID(), epe.event.EventType())
}

type AggregateSaverError struct {
	err       error
	aggregate eventsourcing.Aggregate
}

func (ale AggregateSaverError) Error() string {
	return fmt.Sprintf("failed to save aggregate: %s (aggregate_id %s)",
		ale.err, ale.aggregate.GetAggregateID())
}

type AggregateLoaderError struct {
	err         error
	aggregateID uuid.UUID
}

func (ale AggregateLoaderError) Unwrap() error {
	return ale.err
}

func (ale AggregateLoaderError) Error() string {
	return fmt.Sprintf("failed to load aggregate: %s (aggregate_id %s)",
		ale.err, ale.aggregateID)
}

type TransactionError struct {
	err error
	msg string
}

func (te TransactionError) Error() string {
	return fmt.Sprintf("transaction error %s: %v", te.msg, te.err)
}

func (te TransactionError) Unwrap() error {
	return te.err
}

type EventUnmarshalError struct {
	err        error
	eventModel *EventModel
}

func (eue EventUnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal event: %s (event type %s, aggregate_id %d)",
		eue.err, eue.eventModel.EventType, eue.eventModel.AggregateID)
}

type EventMarshalError struct {
	err   error
	event eventsourcing.Event
}

func (eue EventMarshalError) Error() string {
	return fmt.Sprintf("failed to marshal event: %s (event type %s, aggregate_id %d)",
		eue.err, eue.event.EventType(), eue.event.GetAggregateID())
}

type EventVersionConflictError struct {
	err   error
	event eventsourcing.Event
}

func (cee *EventVersionConflictError) Error() string {
	return fmt.Sprintf("event version conflicted: %s (event type %s, aggregate_id %d)",
		cee.err, cee.event.EventType(), cee.event.GetAggregateID())
}
