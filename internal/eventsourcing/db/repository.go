package db

import (
	"context"
	"reflect"

	"github.com/samwang0723/jarvis/internal/database"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"gorm.io/gorm"
)

type AggregateRepository struct {
	aggregateType reflect.Type

	eventStore EventStore

	projectors map[eventsourcing.EventType][]eventsourcing.Projector

	// aggregateLoader loads aggregates from a projected table.
	aggregateLoader eventsourcing.AggregateLoader

	// aggregateSaver saves aggregates to a projected table.
	aggregateSaver eventsourcing.AggregateSaver
}

type AggregateRepositoryOption func(*AggregateRepository)

// WithAggregateLoader configure a AggregateLoader for the repository.
func WithAggregateLoader(aggregateLoader eventsourcing.AggregateLoader) AggregateRepositoryOption {
	return func(ar *AggregateRepository) {
		ar.aggregateLoader = aggregateLoader
	}
}

// WithAggregateSaver configure a AggregateSaver for the repository.
func WithAggregateSaver(aggregateSaver eventsourcing.AggregateSaver) AggregateRepositoryOption {
	return func(ar *AggregateRepository) {
		ar.aggregateSaver = aggregateSaver
	}
}

func NewAggregateRepository(
	aggregate eventsourcing.Aggregate, db *gorm.DB,
	options ...AggregateRepositoryOption,
) *AggregateRepository {
	eventTable := aggregate.EventTable()

	eventRegistry := eventsourcing.NewEventRegistryFromStateMachine(aggregate)

	repo := &AggregateRepository{}
	repo.aggregateType = reflect.TypeOf(aggregate).Elem()
	repo.eventStore = *NewEventStore(eventTable, eventRegistry, db)
	repo.projectors = make(map[eventsourcing.EventType][]eventsourcing.Projector)

	for _, option := range options {
		option(repo)
	}

	return repo
}

// Load loads a aggregate by aggregateID.
//
// If a AggregateLoader is provided, it loads from the AggregateLoader, otherwise
// it loads from the event store.
func (ar *AggregateRepository) Load(
	ctx context.Context, aggregateID uint64,
) (eventsourcing.Aggregate, error) {
	if ar.aggregateLoader != nil {
		return ar.loadFromAggregateLoader(ctx, aggregateID)
	}

	return ar.loadFromEventStore(ctx, aggregateID)
}

// loadFromEventStore loads a aggregate from AggregateLoader.
func (ar *AggregateRepository) loadFromAggregateLoader(
	ctx context.Context, aggregateID uint64,
) (eventsourcing.Aggregate, error) {
	aggregate, err := ar.aggregateLoader.Load(ctx, aggregateID)
	if err != nil {
		return nil, &AggregateLoaderError{
			err:         err,
			aggregateID: aggregateID,
		}
	}

	return aggregate, nil
}

// loadFromEventStore loads a aggregate from event store.
func (ar *AggregateRepository) loadFromEventStore(
	ctx context.Context, aggregateID uint64,
) (eventsourcing.Aggregate, error) {
	aggregate := ar.NewEmptyAggregate()

	aggregate.SetAggregateID(aggregateID)

	events, err := ar.eventStore.Load(ctx, aggregateID, 0)
	if err != nil {
		return aggregate, err
	}

	if len(events) == 0 {
		return aggregate, &AggregateNotFoundError{
			err:         err,
			aggregateID: aggregateID,
		}
	}

	for _, event := range events {
		if err = aggregate.Apply(event); err != nil {
			return aggregate, &ApplyEventError{err: err, aggregateID: aggregateID}
		}
	}

	return aggregate, nil
}

// Save saves the aggregate.
//
// After appending the uncommitted events to event store, it run projectors.
// If a AggregateSaver is provided, it also save the aggregate with the AggregateSaver.
// If any errors happened, all changes will be rollbacked.
func (ar *AggregateRepository) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
	dbPool := ar.eventStore.db
	if tx, ok := database.GetTx(ctx); ok {
		dbPool = tx
	}

	return dbPool.Transaction(func(_ *gorm.DB) error {
		changes := aggregate.GetChanges()

		if err := ar.eventStore.Append(ctx, changes); err != nil {
			return err
		}

		if ar.aggregateSaver != nil {
			if err := ar.aggregateSaver.Save(ctx, aggregate); err != nil {
				return AggregateSaverError{
					err:       err,
					aggregate: aggregate,
				}
			}
		}

		// run projector
		for _, event := range changes {
			if err := ar.project(ctx, event); err != nil {
				return err
			}
		}

		return nil
	})
}

// project runs projectors.
func (ar *AggregateRepository) project(ctx context.Context, event eventsourcing.Event) error {
	for _, projector := range ar.projectors[event.EventType()] {
		if err := projector.Handle(ctx, event); err != nil {
			return &EventProjectorError{
				err:   err,
				event: event,
			}
		}
	}

	return nil
}

// AddProjector add a projector to repository.
func (ar *AggregateRepository) AddProjector(event eventsourcing.Event, projector eventsourcing.Projector) {
	eventType := event.EventType()
	ar.projectors[eventType] = append(ar.projectors[eventType], projector)
}

func (ar *AggregateRepository) NewEmptyAggregate( //nolint: ireturn // Aggregate is a interface
) eventsourcing.Aggregate {
	//nolint: errcheck // only Aggregate can be registered
	aggregate, _ := reflect.New(ar.aggregateType).Interface().(eventsourcing.Aggregate)

	return aggregate
}
