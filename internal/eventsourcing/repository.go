package eventsourcing

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type AggregateRepository interface {
	Load(ctx context.Context, id entity.ID) (Aggregate, error)
	Save(ctx context.Context, aggregate Aggregate) error
}
