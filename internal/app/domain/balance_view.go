package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

// Define state machine
const (
	balanceInitState    eventsourcing.State = ""
	balanceCreatedState eventsourcing.State = "created"
)

type BalanceView struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	eventsourcing.BaseAggregate
	Balance   float32
	Pending   float32
	Available float32
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

func NewBalanceView(userID uuid.UUID, initBalance float32) (*BalanceView, error) {
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
	bv.AppendChanges(event)

	return bv, nil
}

// MoveAvailableToPending moves balance from available state to pending state.
// It doesn't change total balance.
func (bv *BalanceView) MoveAvailableToPending(transaction *Transaction) error {
	// create a event
	pendingDelta := float32(0.0)
	availableDelta := float32(0.0)
	amount := transaction.CreditAmount - transaction.DebitAmount

	event := &BalanceChanged{
		AvailableDelta: availableDelta - abs(amount),
		PendingDelta:   pendingDelta + abs(amount),
		Amount:         amount,
		Currency:       "TWD",
		TransactionID:  transaction.ID,
		OrderType:      transaction.OrderType,
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
	bv.AppendChanges(event)

	return nil
}

func (bv *BalanceView) MovePendingToAvailable(transaction *Transaction) error {
	pendingDelta := float32(0.0)
	availableDelta := float32(0.0)
	amount := transaction.CreditAmount - transaction.DebitAmount

	event := &BalanceChanged{
		AvailableDelta: availableDelta + abs(amount),
		PendingDelta:   pendingDelta - abs(amount),
		Amount:         amount,
		Currency:       "TWD",
		TransactionID:  transaction.ID,
		OrderType:      transaction.OrderType,
	}

	event.SetAggregateID(bv.GetAggregateID())
	event.SetVersion(bv.Version + 1)
	event.SetCreatedAt(time.Now())

	if err := bv.Apply(event); err != nil {
		return err
	}

	bv.AppendChanges(event)

	return nil
}

func (bv *BalanceView) CreditPending(transaction *Transaction) error {
	pendingDelta := float32(0.0)
	availableDelta := float32(0.0)
	amount := transaction.CreditAmount - transaction.DebitAmount

	event := &BalanceChanged{
		AvailableDelta: availableDelta,
		PendingDelta:   pendingDelta + abs(amount),
		Amount:         amount,
		Currency:       "TWD",
		TransactionID:  transaction.ID,
		OrderType:      transaction.OrderType,
	}

	event.SetAggregateID(bv.GetAggregateID())
	event.SetVersion(bv.Version + 1)
	event.SetCreatedAt(time.Now())

	if err := bv.Apply(event); err != nil {
		return err
	}

	bv.AppendChanges(event)

	return nil
}

func (bv *BalanceView) DebitPending(transaction *Transaction) error {
	pendingDelta := float32(0.0)
	availableDelta := float32(0.0)
	amount := transaction.CreditAmount - transaction.DebitAmount

	event := &BalanceChanged{
		AvailableDelta: availableDelta,
		PendingDelta:   pendingDelta - abs(amount),
		Amount:         amount,
		Currency:       "TWD",
		TransactionID:  transaction.ID,
		OrderType:      transaction.OrderType,
	}

	event.SetAggregateID(bv.GetAggregateID())
	event.SetVersion(bv.Version + 1)
	event.SetCreatedAt(time.Now())

	if err := bv.Apply(event); err != nil {
		return err
	}

	bv.AppendChanges(event)

	return nil
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}

	return x
}
