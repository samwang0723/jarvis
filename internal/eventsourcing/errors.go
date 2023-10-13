package eventsourcing

import "fmt"

type NoTransitionError struct {
	event Event
	state State
}

func (nte *NoTransitionError) Error() string {
	return fmt.Sprintf("no valid transition: event %s, state %s", nte.event.EventType(), nte.state)
}

type EventNotRegisteredError struct {
	eventType EventType
}

func (enre *EventNotRegisteredError) Error() string {
	return fmt.Sprintf("event not registered: %s", enre.eventType)
}
