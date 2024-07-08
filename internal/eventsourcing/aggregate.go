package eventsourcing

import "github.com/gofrs/uuid/v5"

type Aggregate interface {
	Apply(event Event) error
	GetChanges() []Event
	AppendChange(events Event)
	SetAggregateID(aggregateID uuid.UUID)
	GetAggregateID() uuid.UUID
	GetVersion() int
	SetVersion(version int)
	EventTable() string

	StateMachine
}

type BaseAggregate struct {
	ID                uuid.UUID
	Version           int
	uncommittedEvents []Event
}

func (ba *BaseAggregate) AppendChange(event Event) {
	ba.uncommittedEvents = append(ba.uncommittedEvents, event)
}

func (ba *BaseAggregate) GetChanges() []Event {
	return ba.uncommittedEvents
}

func (ba *BaseAggregate) GetVersion() int {
	return ba.Version
}

func (ba *BaseAggregate) SetVersion(version int) {
	ba.Version = version
}

func (ba *BaseAggregate) GetAggregateID() uuid.UUID {
	return ba.ID
}

func (ba *BaseAggregate) SetAggregateID(aggregateID uuid.UUID) {
	ba.ID = aggregateID
}

func (ba *BaseAggregate) GetCurrentState() State {
	return ""
}

func (ba *BaseAggregate) GetTransitions() []Transition {
	return []Transition{}
}

func (ba *BaseAggregate) SkipTransition(_ Event) bool {
	return false
}
