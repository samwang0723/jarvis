package entity

import "github.com/samwang0723/jarvis/internal/eventsourcing"

type BalanceCreated struct {
	eventsourcing.BaseEvent
	InitialBalance float32
}

// EventType returns the name of event
func (*BalanceCreated) EventType() eventsourcing.EventType {
	return "balance.created"
}

type BalanceChanged struct {
	Currency  string
	OrderType string
	eventsourcing.BaseEvent
	TransactionID  uint64
	AvailableDelta float32
	PendingDelta   float32
	Amount         float32
}

func (*BalanceChanged) EventType() eventsourcing.EventType {
	return "balance.changed"
}
