package eventsourcing

type Aggregate interface {
	Apply(event Event) error
	GetChanges() []Event
	AppendChange(events Event)
	SetAggregateID(aggregateID uint64)
	GetAggregateID() uint64
	GetVersion() int
	SetVersion(version int)
	EventTable() string

	StateMachine
}

type BaseAggregate struct {
	ID                uint64 `gorm:"column:id"`
	Version           int    `gorm:"column:version"` // event version number, used for ordering events
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

func (ba *BaseAggregate) GetAggregateID() uint64 {
	return ba.ID
}

func (ba *BaseAggregate) SetAggregateID(aggregateID uint64) {
	ba.ID = aggregateID
}

func (ba *BaseAggregate) GetCurrentState() State {
	return ""
}

func (ba *BaseAggregate) GetTransitions() []Transition {
	return []Transition{}
}

func (ba *BaseAggregate) SkipTransition(event Event) bool {
	return false
}
