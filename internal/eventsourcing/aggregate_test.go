package eventsourcing_test

import (
	"flag"
	"math/rand"
	"os"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
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

func TestAggregate(t *testing.T) {
	t.Parallel()

	aggregate := eventsourcing.BaseAggregate{}

	aggregateID := uuid.Must(uuid.NewV4())
	version := rand.Int() //nolint: gosec // CPRNG is not need for version

	// set and get aggregateID
	aggregate.SetAggregateID(aggregateID)
	assert.Equal(t, aggregateID, aggregate.GetAggregateID())

	// set and get version
	aggregate.SetVersion(version)

	assert.Equal(t, version, aggregate.GetVersion())

	// append and get changes
	event := &testCreatedEvent{}
	aggregate.AppendChanges(event)

	changes := aggregate.GetChanges()
	assert.Equal(t, 1, len(changes))
	assert.Equal(t, event, changes[0])

	// no transitions
	transitions := aggregate.GetTransitions()
	assert.Equal(t, 0, len(transitions))

	// no current state
	state := aggregate.GetCurrentState()
	assert.Equal(t, "", string(state))

	assert.Panics(t, func() {
		aggregate.Apply(event)
	})

	assert.Panics(t, func() {
		aggregate.EventTable()
	})

	assert.False(t, aggregate.SkipTransition(event))
}
