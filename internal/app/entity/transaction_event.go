package entity

import "github.com/samwang0723/jarvis/internal/eventsourcing"

type TransactionCreated struct {
	DebitAmount  float32
	CreditAmount float32
	OrderType    string
	OrderID      uint64

	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*TransactionCreated) EventType() eventsourcing.EventType {
	return "transaction.created"
}

type TransactionCompleted struct {
	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*TransactionCompleted) EventType() eventsourcing.EventType {
	return "transaction.completed"
}

type TransactionFailed struct {
	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*TransactionFailed) EventType() eventsourcing.EventType {
	return "transaction.failed"
}
