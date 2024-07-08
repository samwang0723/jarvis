package db

import (
	"reflect"
	"time"

	"github.com/gofrs/uuid/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

type EventModel struct {
	CreatedAt   time.Time `gorm:"column:created_at"`
	EventType   string    `gorm:"column:event_type"`
	Payload     string    `gorm:"column:payload"`
	AggregateID uuid.UUID `gorm:"column:aggregate_id"`
	ParentID    uuid.UUID `gorm:"column:parent_id"`
	Version     int       `gorm:"column:version"`
}

// ToEvent converts a EventModel to Event.
func (ev *EventModel) ToEvent( //nolint: ireturn // concrete type unknown here
	reg *eventsourcing.EventRegistry,
) (eventsourcing.Event, error) {
	event, err := reg.GetInstance(eventsourcing.EventType(ev.EventType))
	if err != nil {
		return nil, err //nolint: wrapcheck // error is defined in eventsourcing package
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.UnmarshalFromString(ev.Payload, event)
	if err != nil {
		return nil, &EventUnmarshalError{err: err, eventModel: ev}
	}

	event.SetAggregateID(ev.AggregateID)
	event.SetCreatedAt(ev.CreatedAt)
	event.SetVersion(ev.Version)
	event.SetParentID(ev.ParentID)

	return event, nil
}

// NewEventModelFromEvent converts a event to EventModel.
func NewEventModelFromEvent(event eventsourcing.Event) (*EventModel, error) {
	payload := eventPayload(event)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	payloadJSON, err := json.MarshalToString(payload)
	if err != nil {
		return nil, &EventMarshalError{
			err:   err,
			event: event,
		}
	}

	return &EventModel{
		AggregateID: event.GetAggregateID(),
		Version:     event.GetVersion(),
		ParentID:    event.GetParentID(),
		EventType:   string(event.EventType()),
		Payload:     payloadJSON,
		CreatedAt:   event.GetCreatedAt(),
	}, nil
}

// EventPayload return payload of event.
func eventPayload(event eventsourcing.Event) map[string]any {
	payload := make(map[string]any)

	eventType := reflect.TypeOf(event).Elem()
	value := reflect.ValueOf(event).Elem()

	for i := 0; i < eventType.NumField(); i++ {
		field := eventType.Field(i)

		// skip fields in embedded BaseEvent
		if field.Name == "BaseEvent" && field.Type.Kind() == reflect.Struct {
			continue
		}

		payload[field.Name] = value.FieldByName(field.Name).Interface()
	}

	return payload
}
