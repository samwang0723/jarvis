package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/samwang0723/jarvis/internal/eventsourcing"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const listEventsByAggregateIDAndVersion = `
SELECT aggregate_id,
       version,
       parent_id,
       event_type,
       payload,
       created_at
FROM %s
WHERE aggregate_id = $1 and version >= $2
ORDER by version asc
`

const insertEvent = `
INSERT INTO %s (aggregate_id, version, parent_id, event_type, payload, created_at) VALUES ($1, $2, $3, $4, $5, $6)
`

const createEventTable = `
CREATE TABLE %s (
  aggregate_id uuid NOT NULL,
  version int NOT NULL,
  parent_id uuid NOT NULL,
  event_type VARCHAR (50),
  payload jsonb NOT NULL,
  created_at timestamp without time zone NOT NULL,
  PRIMARY KEY (aggregate_id, version)
);
`

type EventStore struct {
	eventTable string
	registry   *eventsourcing.EventRegistry
	dbPool     *pgxpool.Pool
}

func NewEventStore(
	eventTable string, er *eventsourcing.EventRegistry, dbPool *pgxpool.Pool,
) *EventStore {
	store := &EventStore{}
	store.eventTable = eventTable
	store.registry = er
	store.dbPool = dbPool

	return store
}

func (es *EventStore) Load(
	ctx context.Context, aggregateID uuid.UUID, startVersion int,
) ([]eventsourcing.Event, error) {
	events := []eventsourcing.Event{}
	err := Transaction(ctx, es.dbPool, func(ctx context.Context, tx pgx.Tx) error {
		sql := fmt.Sprintf(listEventsByAggregateIDAndVersion, es.eventTable)

		rows, err := tx.Query(ctx, sql, aggregateID, startVersion)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).
				Str("event_table", es.eventTable).
				Str("aggregate_id", aggregateID.String()).
				Msg("load events")

			return fmt.Errorf("load events failed: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var evModel EventModel

			if err := rows.Scan(&evModel.AggregateID,
				&evModel.Version,
				&evModel.ParentID,
				&evModel.EventType,
				&evModel.Payload,
				&evModel.CreatedAt); err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).
					Str("aggregate_id", aggregateID.String()).
					Int("start_version", startVersion).
					Msg("scan model error")

				return fmt.Errorf("scan event model error: %w", err)
			}

			event, err := evModel.ToEvent(es.registry)
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).
					Str("event_type", evModel.EventType).
					Msg("failed converting event model to event")

				return err
			}

			events = append(events, event)
		}

		return nil
	})

	return events, err
}

// Append inserts events to database.
func (es *EventStore) Append(ctx context.Context, events []eventsourcing.Event) error {
	return Transaction(ctx, es.dbPool, func(ctx context.Context, pgtx pgx.Tx) error {
		insertSQL := fmt.Sprintf(insertEvent, es.eventTable)

		for _, event := range events {
			evModel, err := NewEventModelFromEvent(event)
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).Msg("failed to marshal event")

				return err
			}

			_, err = pgtx.Exec(ctx, insertSQL,
				evModel.AggregateID,
				evModel.Version,
				evModel.ParentID,
				evModel.EventType,
				evModel.Payload,
				evModel.CreatedAt)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
					return &EventVersionConflictError{
						err:   err,
						event: event,
					}
				}

				zerolog.Ctx(ctx).Warn().Err(err).
					Str("table", es.eventTable).
					Str("sql", insertSQL).
					Str("aggregate_id", event.GetAggregateID().String()).
					Str("event_type", string(event.EventType())).
					Msg("insert event error")

				return fmt.Errorf("insert event %t error %w", event, err)
			}
		}

		return nil
	})
}

// Migration returns a sql for creating the event table.
func (es *EventStore) Migration() string {
	return fmt.Sprintf(createEventTable, es.eventTable)
}
