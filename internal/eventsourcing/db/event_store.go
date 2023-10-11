package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"gorm.io/gorm"
)

const listEventsByAggregateIDAndVersion = `
SELECT aggregate_id,
       version,
       parent_id,
       event_type,
       payload,
       created_at
FROM %s
WHERE aggregate_id = ? and version >= ?
ORDER by version asc
`

const insertEvent = `
INSERT INTO %s (aggregate_id, version, parent_id, event_type, payload, created_at) 
        VALUES (?, ?, ?, ?, ?, ?)
`

type EventStore struct {
	eventTable string
	registry   *eventsourcing.EventRegistry
	db         *gorm.DB
}

func NewEventStore(
	eventTable string, er *eventsourcing.EventRegistry, db *gorm.DB,
) *EventStore {
	store := &EventStore{}
	store.eventTable = eventTable
	store.registry = er
	store.db = db

	return store
}

func (es *EventStore) Load(
	ctx context.Context, aggregateID uint64, startVersion int,
) ([]eventsourcing.Event, error) {
	events := []eventsourcing.Event{}
	err := es.db.Transaction(func(tx *gorm.DB) error {
		sql := fmt.Sprintf(listEventsByAggregateIDAndVersion, es.eventTable)
		rows, err := tx.Raw(sql, aggregateID, startVersion).Rows()
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).
				Str("event_table", es.eventTable).
				Uint64("aggregate_id", aggregateID).
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
					Uint64("aggregate_id", aggregateID).
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

func (es *EventStore) Append(ctx context.Context, events []eventsourcing.Event) error {
	return es.db.Transaction(func(tx *gorm.DB) error {
		for _, event := range events {
			evModel, err := NewEventModelFromEvent(event)
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).
					Str("event_type", string(event.EventType())).
					Msg("failed converting event to event model")

				return err
			}

			sql := fmt.Sprintf(insertEvent, es.eventTable)
			if err := tx.Exec(sql,
				evModel.AggregateID,
				evModel.Version,
				evModel.ParentID,
				evModel.EventType,
				evModel.Payload,
				evModel.CreatedAt).Error; err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).
					Str("event_type", string(event.EventType())).
					Msg("failed inserting event")

				return err
			}
		}

		return nil
	})
}
