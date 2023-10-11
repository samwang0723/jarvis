package eventsourcing

import (
	"context"
)

type AggregateLoader interface {
	Load(ctx context.Context, id uint64) (Aggregate, error)
}
