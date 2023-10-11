// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package entity

import (
	"time"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

// Define state machine
const (
	balanceInitState    eventsourcing.State = ""
	balanceCreatedState eventsourcing.State = "created"
)

type BalanceView struct {
	Balance   float32   `gorm:"column:balance" json:"balance"`
	Pending   float32   `gorm:"column:pending" json:"pending"`
	Available float32   `gorm:"column:available" json:"available"`
	CreatedAt time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" mapstructure:"updated_at"`

	eventsourcing.BaseAggregate
}

func (BalanceView) TableName() string {
	return "balance_views"
}

func (bv *BalanceView) EventTable() string {
	return "balance_events"
}

// Apply updates the aggregate according to a event.
func (bv *BalanceView) Apply(event eventsourcing.Event) error {
	switch event := event.(type) {
	case *BalanceCreated:
		bv.Balance = event.InitialBalance
		bv.Available = event.InitialBalance

		bv.CreatedAt = event.CreatedAt
		bv.UpdatedAt = event.CreatedAt
		bv.SetAggregateID(event.AggregateID)

	case *BalanceChanged:
		bv.Available += event.AvailableDelta
		bv.Pending += event.PendingDelta
		bv.Balance += event.AvailableDelta + event.PendingDelta
	default:
		return &UnsupportedEventError{event: event}
	}

	bv.Version = event.GetVersion()

	return nil
}

// GetStates returns all possible state transitions
func (bv *BalanceView) GetTransitions() []eventsourcing.Transition {
	return []eventsourcing.Transition{
		{
			FromState: balanceInitState,
			Event:     &BalanceCreated{},
			ToState:   balanceCreatedState,
		},
		{
			FromState: balanceCreatedState,
			Event:     &BalanceChanged{},
			ToState:   balanceCreatedState,
		},
	}
}

func NewBalanceView(userID uint64, initBalance float32) (*BalanceView, error) {
	// create a init balance_view
	bv := &BalanceView{}

	// create a event
	event := &BalanceCreated{
		InitialBalance: initBalance,
	}

	// fill base event data
	event.SetAggregateID(userID)
	event.SetVersion(1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := bv.Apply(event); err != nil {
		return nil, err
	}
	// record uncommitted events
	bv.AppendChange(event)

	return bv, nil
}

// MoveAvailableToPending moves balance from available state to pending state.
// It doesn't change total balance.
func (bv *BalanceView) MoveAvailableToPending(amount float32) error {
	// create a event
	pendingDelta := float32(0.0)
	availableDelta := float32(0.0)

	event := &BalanceChanged{
		AvailableDelta: availableDelta - amount,
		PendingDelta:   pendingDelta + amount,
	}
	// fill base event data
	event.SetAggregateID(bv.GetAggregateID())
	event.SetVersion(bv.Version + 1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := bv.Apply(event); err != nil {
		return err
	}
	// record uncommitted events
	bv.AppendChange(event)

	return nil
}
