package eventsourcing

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type AggregateLoader interface {
	Load(ctx context.Context, id uuid.UUID) (Aggregate, error)
}
