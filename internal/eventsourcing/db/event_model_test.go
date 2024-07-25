package db_test

import (
	"flag"
	"os"
	"testing"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/eventsourcing/db"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

type payloadTestEvent struct {
	Decimal *decimal.Big
	Str     string
	eventsourcing.BaseEvent
	Int  int
	UUID uuid.UUID
}

func (pte *payloadTestEvent) EventType() eventsourcing.EventType {
	return "payload_test"
}

func TestEventModel(t *testing.T) {
	t.Parallel()

	event := &payloadTestEvent{
		Int:     1,
		UUID:    uuid.Must(uuid.NewV4()),
		Decimal: decimal.New(1234, 2),
		Str:     "str",
	}
	registry := eventsourcing.NewEventRegistry()

	// convert a event to eventModel
	eventModel, err := db.NewEventModelFromEvent(event)
	assert.Nil(t, err)

	// convert a eventModel back to event, failed without event registry
	_, err = eventModel.ToEvent(registry)
	assert.NotNil(t, err)

	expectedErr := &eventsourcing.EventNotRegisteredError{}
	assert.ErrorAs(t, err, &expectedErr)

	// register event
	registry.Register(&payloadTestEvent{})

	event2, err := eventModel.ToEvent(registry)
	assert.Nil(t, err)
	assert.Equal(t, event, event2)
}
