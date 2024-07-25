package eventsourcing

type State string

type Transition struct {
	Event     Event
	FromState State
	ToState   State
}

type StateMachine interface {
	// GetStates returns current state of a aggregate.
	GetCurrentState() State

	// GetTransitions return all possible state transition for a aggregate
	GetTransitions() []Transition

	// SkipTransition returns true if the transition should be skipped.
	SkipTransition(event Event) bool
}

// TransistOnEvent returns the new state given current state and event.
func TransistOnEvent(stateMachine StateMachine, event Event) (State, error) {
	if err := event.Validate(); err != nil {
		return "", &InvalidEventError{
			Type: event.EventType(),
			Err:  err,
		}
	}

	currentState := stateMachine.GetCurrentState()

	if stateMachine.SkipTransition(event) {
		return currentState, nil
	}

	transitions := stateMachine.GetTransitions()
	eventType := event.EventType()

	for _, transition := range transitions {
		if transition.FromState == currentState && transition.Event.EventType() == eventType {
			return transition.ToState, nil
		}
	}

	return "", &NoTransitionError{
		event: event,
		state: currentState,
	}
}
