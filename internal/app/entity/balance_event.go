package entity

import "github.com/samwang0723/jarvis/internal/eventsourcing"

type BalanceCreated struct {
	InitialBalance float32

	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*BalanceCreated) EventType() eventsourcing.EventType {
	return "balance.created"
}

type BalanceChanged struct {
	AvailableDelta float32
	PendingDelta   float32

	eventsourcing.BaseEvent
}

func (*BalanceChanged) EventType() eventsourcing.EventType {
	return "balance.changed"
}
