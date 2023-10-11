package eventsourcing

import "context"

type EventStore[T Aggregate] interface {
	Load(ctx context.Context, aggregateID uint64, startVersion int) ([]Event, error)
	Append(ctx context.Context, events []Event) error
}
