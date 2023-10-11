package eventsourcing

import "context"

type AggregateSaver interface {
	Save(ctx context.Context, aggregate Aggregate) error
}
