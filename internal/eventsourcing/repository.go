package eventsourcing

import (
	"context"
)

type AggregateRepository interface {
	Load(ctx context.Context, id uint64) (Aggregate, error)
	Save(ctx context.Context, aggregate Aggregate) error
}
