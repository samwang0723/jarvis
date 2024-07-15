package domain

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

type DataValidationError struct {
	dataType string
}

func (e *DataValidationError) Error() string {
	return fmt.Sprintf("validation failed: invalid %s", e.dataType)
}

type DataMissingError struct {
	dataType string
}

func (e *DataMissingError) Error() string {
	return fmt.Sprintf("data missing: %s", e.dataType)
}
