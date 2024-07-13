package eventsourcing

import "reflect"

type EventRegistry struct {
	events map[EventType]reflect.Type
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[EventType]reflect.Type),
	}
}

func NewEventRegistryFromStateMachine(stateMachine StateMachine) *EventRegistry {
	registry := NewEventRegistry()

	for _, transition := range stateMachine.GetTransitions() {
		registry.Register(transition.Event)
	}

	return registry
}

func (er *EventRegistry) Register(event Event) {
	er.events[event.EventType()] = reflect.TypeOf(event).Elem()
}

func (er *EventRegistry) Get(eventType EventType) (reflect.Type, error) {
	typ, ok := er.events[eventType]
	if !ok {
		return nil, &EventNotRegisteredError{event: eventType}
	}

	return typ, nil
}

func (er *EventRegistry) GetInstance(eventType EventType) (Event, error) {
	typ, err := er.Get(eventType)
	if err != nil {
		return nil, err
	}

	//nolint: errcheck // only Event can be registered
	event, _ := reflect.New(typ).Interface().(Event)

	return event, nil
}
