package db_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/common/remotetest"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/eventsourcing/db"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// create database
	pool := remotetest.SetupPostgresClient(t, true)

	registry := eventsourcing.NewEventRegistryFromStateMachine(&testAggregate{})

	repository := db.NewAggregateRepository(&testAggregate{}, pool)
	eventStore := db.NewEventStore((&testAggregate{}).EventTable(), registry, pool)
	aggregateID := uuid.Must(uuid.NewV4())

	// create event table
	_, err := pool.Exec(ctx, eventStore.Migration())
	assert.Nil(t, err)

	// aggregate does exist at the beginning
	_, err = repository.Load(ctx, aggregateID)
	assert.NotNil(t, err)

	// create aggregate
	aggregate := &testAggregate{}
	aggregate.SetAggregateID(aggregateID)

	// apply an event
	event := &testedEvent{}
	event.SetAggregateID(aggregateID)
	event.SetVersion(1)
	aggregate.Apply(event)
	aggregate.AppendChanges(event) // uncommitted events

	// save aggregate
	err = repository.Save(ctx, aggregate) // commit aggregate and events, may be rejected
	assert.Nil(t, err)

	// load events
	events, err := eventStore.Load(ctx, aggregateID, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))

	// load aggregate from db
	loadedAgg, err := repository.Load(ctx, aggregateID)
	assert.Nil(t, err)
	assert.Equal(t, 1, loadedAgg.GetVersion())
}

type testAggregateLoader struct {
	fail bool
}

// nolint: ireturn // follow AggregateLoader interface
func (tal *testAggregateLoader) Load(
	context.Context, uuid.UUID,
) (eventsourcing.Aggregate, error) {
	if tal.fail {
		return nil, db.AggregateLoaderError{}
	}

	return &testAggregate{}, nil
}

func TestRepository_AggregateLoader(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pool := remotetest.SetupPostgresClient(t, true)

	registry := eventsourcing.NewEventRegistry()
	registry.Register(&testedEvent{})

	tests := []struct {
		expectErr  any
		repository *db.AggregateRepository
		name       string
	}{
		{
			name: "happy path: aggregate loader success",
			repository: db.NewAggregateRepository(&testAggregate{}, pool,
				db.WithAggregateLoader(&testAggregateLoader{fail: false})),
			expectErr: nil,
		},
		{
			name: "unhappy path: aggregate loader failed",
			repository: db.NewAggregateRepository(&testAggregate{}, pool,
				db.WithAggregateLoader(&testAggregateLoader{fail: true})),
			expectErr: &db.AggregateLoaderError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			aggregateID := uuid.Must(uuid.NewV4())

			_, err := tt.repository.Load(ctx, aggregateID)
			if tt.expectErr != nil {
				assert.NotNil(t, err)
				assert.ErrorAs(t, err, &tt.expectErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
