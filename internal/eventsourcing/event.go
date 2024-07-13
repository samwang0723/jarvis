package eventsourcing

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type EventType string

type Event interface {
	// EventType returns the name of event.
	EventType() EventType

	// GetAggregateID returns event's aggregate id.
	GetAggregateID() uuid.UUID

	// SetAggregateID changes event's aggregate id.
	SetAggregateID(uuid.UUID)

	// GetParentID returns event's parent id
	GetParentID() uuid.UUID

	// SetParentID changes event's parent id
	SetParentID(uuid.UUID)

	// GetVersion returns event's version
	GetVersion() int

	// SetVersion changes event's version
	SetVersion(int)

	// GetCreatedAt returns event's create time
	GetCreatedAt() time.Time

	// SetCreatedAt changes event's create time
	SetCreatedAt(time.Time)

	// Validate checks if event is valid
	Validate() error
}

// BaseEvent implements common functionality for Event interface.
//
// It doesn't implement EventType method.
type BaseEvent struct {
	CreatedAt   time.Time
	Version     int
	ID          uuid.UUID
	AggregateID uuid.UUID
	ParentID    uuid.UUID
}

var _ Event = (*BaseEvent)(nil)

func (*BaseEvent) EventType() EventType {
	panic("implement me in derived event")
}

func (be *BaseEvent) GetAggregateID() uuid.UUID {
	return be.AggregateID
}

func (be *BaseEvent) SetAggregateID(id uuid.UUID) {
	be.AggregateID = id
}

func (be *BaseEvent) GetParentID() uuid.UUID {
	return be.ParentID
}

func (be *BaseEvent) SetParentID(parentID uuid.UUID) {
	be.ParentID = parentID
}

func (be *BaseEvent) GetVersion() int {
	return be.Version
}

func (be *BaseEvent) SetVersion(version int) {
	be.Version = version
}

func (be *BaseEvent) GetCreatedAt() time.Time {
	return be.CreatedAt
}

func (be *BaseEvent) SetCreatedAt(createdAt time.Time) {
	be.CreatedAt = createdAt
}

func (be *BaseEvent) Validate() error {
	return nil // no validation by default
}
