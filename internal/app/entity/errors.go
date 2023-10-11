package entity

import (
	"fmt"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

type UnsupportedEventError struct {
	event eventsourcing.Event
}

func (e *UnsupportedEventError) Error() string {
	return fmt.Sprintf("unsupported event: %v", e.event.EventType())
}
