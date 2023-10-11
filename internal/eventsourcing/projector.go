package eventsourcing

import "context"

type Projector interface {
	Handle(ctx context.Context, event Event) error
}
