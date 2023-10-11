package eventsourcing

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type AggregateLoader interface {
	Load(ctx context.Context, id entity.ID) (Aggregate, error)
}
