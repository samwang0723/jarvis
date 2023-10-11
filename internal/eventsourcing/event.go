package eventsourcing

import (
	"time"

	"github.com/samwang0723/jarvis/internal/app/entity"
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
	GetAggregateId() uint64
	SetAggregateId(uint64)
	GetParentId() uint64
	SetParentId(uint64)
	GetVersion() int
	SetVersion(int)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
}

type BaseEvent struct {
	ID          entity.ID
	CreatedAt   time.Time
	AggregateID uint64
	ParentID    uint64
	Version     int
}

func (be *BaseEvent) GetAggregateId() uint64 {
	return be.AggregateID
}

func (be *BaseEvent) SetAggregateId(id uint64) {
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
