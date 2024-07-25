package eventsourcing_test

import (
	"errors"
	"testing"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStateMachine struct {
	state string
}

func (tsm *testStateMachine) GetCurrentState() eventsourcing.State {
	return eventsourcing.State(tsm.state)
}

type testCreatedEvent struct {
	eventsourcing.BaseEvent
}

func (tce *testCreatedEvent) EventType() eventsourcing.EventType {
	return "test.created"
}

type testCreatedEventWithValidator struct {
	eventsourcing.BaseEvent
	validationPassed bool
}

func (tce *testCreatedEventWithValidator) EventType() eventsourcing.EventType {
	return "test.created.with.validator"
}

var errValidationFailed = errors.New("validation failed")

func (tce *testCreatedEventWithValidator) Validate() error {
	if tce.validationPassed {
		return nil
	}

	return errValidationFailed
}

type testInvalidEvent struct {
	eventsourcing.BaseEvent
}

func (tce *testInvalidEvent) EventType() eventsourcing.EventType {
	return "test.invalid"
}

type testSkipTransitionEvent struct {
	eventsourcing.BaseEvent
}

func (tske *testSkipTransitionEvent) EventType() eventsourcing.EventType {
	return "test.skip_transition"
}

func (tsm *testStateMachine) GetTransitions() []eventsourcing.Transition {
	return []eventsourcing.Transition{
		{
			FromState: "init",
			ToState:   "created",
			Event:     &testCreatedEvent{},
		},
		{
			FromState: "init",
			ToState:   "created",
			Event:     &testCreatedEventWithValidator{},
		},
	}
}

func (tsm *testStateMachine) SkipTransition(event eventsourcing.Event) bool {
	return event.EventType() == "test.skip_transition"
}

func TestTransistOnEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectErr   any
		name        string
		state       string
		expectState eventsourcing.State
		events      []eventsourcing.Event
	}{
		{
			name:        "happy path: transition exist",
			state:       "init",
			events:      []eventsourcing.Event{&testCreatedEvent{}},
			expectErr:   nil,
			expectState: "created",
		},
		{
			name:        "unhappy path: transition not exist",
			state:       "init",
			events:      []eventsourcing.Event{&testInvalidEvent{}},
			expectErr:   eventsourcing.NoTransitionError{},
			expectState: "",
		},
		{
			name:        "happy path: skip stata transition",
			state:       "init",
			events:      []eventsourcing.Event{&testCreatedEvent{}, &testSkipTransitionEvent{}},
			expectErr:   nil,
			expectState: "created",
		},
		{
			name:  "happy path: event validation ok",
			state: "init",
			events: []eventsourcing.Event{
				&testCreatedEventWithValidator{validationPassed: true},
			},
			expectErr:   nil,
			expectState: "created",
		},
		{
			name:  "unhappy path: event validation failed",
			state: "init",
			events: []eventsourcing.Event{
				&testCreatedEventWithValidator{validationPassed: false},
			},
			expectErr:   eventsourcing.InvalidEventError{},
			expectState: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stateMachine := testStateMachine{state: tt.state}

			for _, event := range tt.events {
				newState, err := eventsourcing.TransistOnEvent(&stateMachine, event)
				if err != nil && tt.expectErr != nil {
					targetErr := tt.expectErr
					assert.ErrorAs(t, err, &targetErr)
				} else {
					require.Nil(t, err)
					require.Nil(t, tt.expectErr)
				}

				stateMachine.state = string(newState)
			}

			if stateMachine.GetCurrentState() != tt.expectState {
				t.Fatalf("expect state %s, got %s", tt.expectState, stateMachine.GetCurrentState())
			}
		})
	}
}
