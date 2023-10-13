package eventsourcing

type State string

type Transition struct {
	FromState State
	ToState   State
	Event     Event
}

type StateMachine interface {
	GetCurrentState() State
	GetTransitions() []Transition
	SkipTransition(event Event) bool
}

func TransistOnEvent(stateMachine StateMachine, event Event) (State, error) {
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

	return "", &NoTransitionError{state: currentState, event: event}
}
