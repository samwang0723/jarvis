package domain

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

const (
	OrderTypeBuy      = "Buy"
	OrderTypeSell     = "Sell"
	OrderTypeFee      = "Fee"
	OrderTypeTax      = "Tax"
	OrderTypeDeposit  = "Deposit"
	OrderTypeWithdraw = "Withdraw"
)

// Define state machine
const (
	transactionInitState      eventsourcing.State = ""
	transactionCreatedState   eventsourcing.State = "created"
	transactionCompletedState eventsourcing.State = "completed"
	transactionFailedState    eventsourcing.State = "failed"
)

type Transaction struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderType string
	Status    string
	eventsourcing.BaseAggregate
	UserID       uuid.UUID
	OrderID      uuid.UUID
	CreditAmount float32
	DebitAmount  float32
}

// ensure Transaction implements Aggregate interface
var _ eventsourcing.Aggregate = &Transaction{}

func (Transaction) TableName() string {
	return "transactions"
}

func (tran *Transaction) EventTable() string {
	return "transaction_events"
}

func (tran *Transaction) GetCurrentState() eventsourcing.State {
	return eventsourcing.State(tran.Status)
}

// Apply updates the aggregate according to a event.
func (tran *Transaction) Apply(event eventsourcing.Event) error {
	newState, err := eventsourcing.TransistOnEvent(tran, event)
	if err != nil {
		return fmt.Errorf("failed to transition state: %w", err)
	}

	tran.Status = string(newState)

	switch event := event.(type) {
	case *TransactionCreated:
		tran.UserID = event.GetParentID()
		tran.OrderType = event.OrderType
		tran.OrderID = event.OrderID
		tran.CreditAmount = event.CreditAmount
		tran.DebitAmount = event.DebitAmount
		tran.CreatedAt = event.CreatedAt
		tran.UpdatedAt = event.CreatedAt
		tran.SetAggregateID(event.AggregateID)
	case *TransactionCompleted:
		tran.UpdatedAt = event.CreatedAt
	case *TransactionFailed:
		tran.UpdatedAt = event.CreatedAt
	default:
		return &UnsupportedEventError{event: event}
	}

	tran.Version = event.GetVersion()

	return nil
}

// GetStates returns all possible state transitions
func (tran *Transaction) GetTransitions() []eventsourcing.Transition {
	return []eventsourcing.Transition{
		{
			FromState: transactionInitState,
			Event:     &TransactionCreated{},
			ToState:   transactionCreatedState,
		},
		{
			FromState: transactionCreatedState,
			Event:     &TransactionCompleted{},
			ToState:   transactionCompletedState,
		},
		{
			FromState: transactionCreatedState,
			Event:     &TransactionFailed{},
			ToState:   transactionFailedState,
		},
	}
}

func NewTransaction(
	userID uuid.UUID,
	orderType string,
	creditAmount float32,
	debitAmount float32,
	orderID ...uuid.UUID,
) (*Transaction, error) {
	id := uuid.Must(uuid.NewV4())
	tran := &Transaction{
		BaseAggregate: eventsourcing.BaseAggregate{
			ID: id,
		},
	}
	event := &TransactionCreated{
		CreditAmount: creditAmount,
		DebitAmount:  debitAmount,
		OrderType:    orderType,
	}

	if len(orderID) > 0 {
		event.OrderID = orderID[0]
	}

	// fill base event data
	event.SetAggregateID(id)
	event.SetParentID(userID)
	event.SetVersion(1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := tran.Apply(event); err != nil {
		return nil, err
	}
	// record uncommitted events
	tran.AppendChange(event)

	return tran, nil
}

func (tran *Transaction) Complete() error {
	event := &TransactionCompleted{}

	// fill base event data
	event.SetAggregateID(tran.ID)
	event.SetParentID(tran.UserID)
	event.SetVersion(tran.Version + 1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := tran.Apply(event); err != nil {
		return err
	}
	// record uncommitted events
	tran.AppendChange(event)

	return nil
}

func (tran *Transaction) Fail() error {
	event := &TransactionFailed{}

	// fill base event data
	event.SetAggregateID(tran.ID)
	event.SetParentID(tran.UserID)
	event.SetVersion(tran.Version + 1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := tran.Apply(event); err != nil {
		return err
	}
	// record uncommitted events
	tran.AppendChange(event)

	return nil
}
