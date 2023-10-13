package eventsourcing

import (
	"time"
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
	GetAggregateID() uint64
	SetAggregateID(uint64)
	GetParentID() uint64
	SetParentID(uint64)
	GetVersion() int
	SetVersion(int)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
}

type BaseEvent struct {
	ID          uint64
	CreatedAt   time.Time
	AggregateID uint64
	ParentID    uint64
	Version     int
}

func (be *BaseEvent) GetAggregateID() uint64 {
	return be.AggregateID
}

func (be *BaseEvent) SetAggregateID(id uint64) {
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

func (be *BaseEvent) GetParentID() uint64 {
	return be.ParentID
}

func (be *BaseEvent) SetParentID(parentID uint64) {
	be.ParentID = parentID
}
