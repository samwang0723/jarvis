package eventsourcing

import "github.com/gofrs/uuid/v5"

type Aggregate interface {
	// Appply updates aggregate based on event.
	Apply(event Event) error

	// Changes returns a list of uncommitted events.
	GetChanges() []Event

	// AppendChanges add a event to a uncommitted event list.
	AppendChanges(event Event)

	// SetAggregateID assigns identifier to aggregate.
	SetAggregateID(uuid.UUID)

	// GetAggregateID returns unique identifier of aggregate.
	GetAggregateID() uuid.UUID

	// GetVersion return current version of the aggregate.
	GetVersion() int

	// SetVersion sets current version of the aggregate.
	SetVersion(int)

	// EventTable return the table name.
	EventTable() string

	StateMachine
}

// BaseAggregate implements common functionality for Aggregate interface.
//
// It doesn't implement Apply, EventTable.
type BaseAggregate struct {
	uncommittedEvents []Event
	Version           int
	ID                uuid.UUID
}

var _ Aggregate = (*BaseAggregate)(nil)

func (ba *BaseAggregate) GetAggregateID() uuid.UUID {
	return ba.ID
}

func (*BaseAggregate) Apply(_ Event) error {
	panic("implement me in derived aggregate")
}

func (*BaseAggregate) EventTable() string {
	panic("implement me in derived aggregate")
}

func (ba *BaseAggregate) SetAggregateID(id uuid.UUID) {
	ba.ID = id
}

func (ba *BaseAggregate) GetChanges() []Event {
	return ba.uncommittedEvents
}

func (ba *BaseAggregate) AppendChanges(event Event) {
	ba.uncommittedEvents = append(ba.uncommittedEvents, event)
}

func (ba *BaseAggregate) GetVersion() int {
	return ba.Version
}

func (ba *BaseAggregate) SetVersion(version int) {
	ba.Version = version
}

func (ba *BaseAggregate) GetCurrentState() State {
	return ""
}

func (ba *BaseAggregate) GetTransitions() []Transition {
	return []Transition{}
}

func (ba *BaseAggregate) SkipTransition(Event) bool {
	return false
}
