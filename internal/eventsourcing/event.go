package eventsourcing

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	EventTransactionPending    = "TransactionPending"
	EventTransactionProcessing = "TransactionProcessing"
	EventTransactionCompleted  = "TransactionCompleted"
	EventTransactionFailed     = "TransactionFailed"
	EventTransactionCancelled  = "TransactionCancelled"
)

type EventType string

type Event interface {
	EventType() EventType
	GetAggregateID() uuid.UUID
	SetAggregateID(uuid.UUID)
	GetParentID() uuid.UUID
	SetParentID(uuid.UUID)
	GetVersion() int
	SetVersion(int)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
}

type BaseEvent struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	AggregateID uuid.UUID
	ParentID    uuid.UUID
	Version     int
}

func (be *BaseEvent) GetAggregateID() uuid.UUID {
	return be.AggregateID
}

func (be *BaseEvent) SetAggregateID(id uuid.UUID) {
	be.AggregateID = id
}

func (be *BaseEvent) GetVersion() int {
	return be.Version
}

func (be *BaseEvent) SetVersion(version int) {
	be.Version = version
}

func (be *BaseEvent) GetCreatedAt() time.Time {
	return be.CreatedAt
}

func (be *BaseEvent) SetCreatedAt(createdAt time.Time) {
	be.CreatedAt = createdAt
}

func (be *BaseEvent) GetParentID() uuid.UUID {
	return be.ParentID
}

func (be *BaseEvent) SetParentID(parentID uuid.UUID) {
	be.ParentID = parentID
}
