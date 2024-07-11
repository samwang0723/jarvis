package db_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/samwang0723/jarvis/internal/common/remotetest"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/eventsourcing/db"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

type testAggregate struct {
	state string
	eventsourcing.BaseAggregate
}

func (ta *testAggregate) EventTable() string {
	return "test_event"
}

func (ta *testAggregate) Apply(event eventsourcing.Event) error {
	ta.Version = event.GetVersion()

	return nil
}

func (ta *testAggregate) GetTransitions() []eventsourcing.Transition {
	return []eventsourcing.Transition{
		{
			FromState: "",
			ToState:   "tested",
			Event:     &testedEvent{},
		},
	}
}

func (ta *testAggregate) GetCurrentState() eventsourcing.State {
	return eventsourcing.State(ta.state)
}

type testedEvent struct {
	eventsourcing.BaseEvent
}

func (te testedEvent) EventType() eventsourcing.EventType {
	return "tested"
}

func TestEventStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool := remotetest.SetupPostgresClient(t, false)

	reg := eventsourcing.NewEventRegistry()
	reg.Register(&testedEvent{})
	eventStore := db.NewEventStore("test_event", reg, pool)

	// migration
	_, err := pool.Exec(ctx, eventStore.Migration())
	assert.Nil(t, err)

	// no events at the beginning
	aggregateID := uuid.Must(uuid.NewV4())
	events, err := eventStore.Load(ctx, aggregateID, 0)
	assert.Nil(t, err)
	assert.Len(t, events, 0)

	// append an event
	event := &testedEvent{}
	event.SetVersion(1)
	event.SetAggregateID(aggregateID)
	err = eventStore.Append(ctx, []eventsourcing.Event{event})
	assert.Nil(t, err)

	// can not insert event of the same version
	err = eventStore.Append(ctx, []eventsourcing.Event{event})
	assert.NotNil(t, err)

	// load event
	events, err = eventStore.Load(ctx, aggregateID, 0)
	assert.Nil(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, reflect.TypeOf(events[0]), reflect.TypeOf(event))
}

func TestEventStore_TableNotExist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool := remotetest.SetupPostgresClient(t, false)

	reg := eventsourcing.NewEventRegistry()
	reg.Register(&testedEvent{})

	eventStore := db.NewEventStore("not_exist_table", reg, pool)

	aggregateID := uuid.Must(uuid.NewV4())

	_, err := eventStore.Load(ctx, aggregateID, 0)
	assert.NotNil(t, err)

	err = eventStore.Append(ctx, []eventsourcing.Event{&testedEvent{}})
	assert.NotNil(t, err)
}

func TestEventStore_EventNotRegistered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool := remotetest.SetupPostgresClient(t, false)

	reg := eventsourcing.NewEventRegistry()

	eventStore := db.NewEventStore("test_event", reg, pool)

	// migration
	_, err := pool.Exec(ctx, eventStore.Migration())
	assert.Nil(t, err)

	aggregateID := uuid.Must(uuid.NewV4())
	event := &testedEvent{}
	event.SetVersion(1)
	event.SetAggregateID(aggregateID)

	err = eventStore.Append(ctx, []eventsourcing.Event{event})
	assert.Nil(t, err)

	_, err = eventStore.Load(ctx, aggregateID, 0)
	assert.NotNil(t, err)
}

func TestEventStore_EventVersionConflict(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool := remotetest.SetupPostgresClient(t, false)

	reg := eventsourcing.NewEventRegistry()

	eventStore := db.NewEventStore("test_event", reg, pool)

	// migration
	_, err := pool.Exec(ctx, eventStore.Migration())
	assert.Nil(t, err)

	aggregateID := uuid.Must(uuid.NewV4())
	event := &testedEvent{}
	event.SetVersion(1)
	event.SetAggregateID(aggregateID)

	// first event success
	err = eventStore.Append(ctx, []eventsourcing.Event{event})
	assert.Nil(t, err)

	// second event with the same version, fail
	err = eventStore.Append(ctx, []eventsourcing.Event{event})

	conflictErr := &db.EventVersionConflictError{}

	assert.ErrorAs(t, err, &conflictErr)
}
