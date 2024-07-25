package eventsourcing_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/stretchr/testify/assert"
)

func TestBaseEvent(t *testing.T) {
	t.Parallel()

	event := eventsourcing.BaseEvent{}

	aggregateID := uuid.Must(uuid.NewV4())
	parentID := uuid.Must(uuid.NewV4())
	version := rand.Int() //nolint: gosec // CPRNG is not need for version
	createdAt := time.Now()

	// set and get aggregateID
	event.SetAggregateID(aggregateID)
	assert.Equal(t, aggregateID, event.GetAggregateID())

	// set and get parentID
	event.SetParentID(parentID)

	assert.Equal(t, parentID, event.GetParentID())

	// set and get version
	event.SetVersion(version)
	assert.Equal(t, version, event.GetVersion())

	// set and get createAt
	event.SetCreatedAt(createdAt)
	assert.Equal(t, createdAt, event.GetCreatedAt())

	assert.Panics(t, func() {
		event.EventType()
	})
}
