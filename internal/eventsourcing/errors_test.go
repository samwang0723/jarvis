package eventsourcing

import (
	"testing"
)

type testTestedEvent struct {
	BaseEvent
}

func (tce *testTestedEvent) EventType() EventType {
	return "test.tested"
}

func TestNoTransitionError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		event Event
		state State
		want  string
	}{
		{
			name:  "happy path",
			event: &testTestedEvent{},
			state: "init",
			want:  "no valid transition: state init, event test.tested",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &NoTransitionError{
				event: tt.event,
				state: tt.state,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("NoTransitionError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventNotRegisteredError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		event EventType
		want  string
	}{
		{
			name:  "happy path",
			event: "created",
			want:  "event not registered: created",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &EventNotRegisteredError{
				event: tt.event,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("EventNotRegisteredError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
