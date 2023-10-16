package entity

import (
	"fmt"
	"time"

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
	UserID       uint64    `gorm:"column:user_id" json:"userId"`
	OrderType    string    `gorm:"column:order_type" json:"orderType"`
	CreditAmount float32   `gorm:"column:credit_amount" json:"creditAmount"`
	DebitAmount  float32   `gorm:"column:debit_amount" json:"debitAmount"`
	OrderID      uint64    `gorm:"column:order_id" json:"orderId"`
	Status       string    `gorm:"column:status" json:"status,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" mapstructure:"updated_at"`

	eventsourcing.BaseAggregate
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
	userID uint64,
	orderType string,
	creditAmount float32,
	debitAmount float32,
	orderID ...uint64,
) (*Transaction, error) {
	id, err := GenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}

	tran := &Transaction{
		UserID:    userID,
		OrderType: orderType,
	}

	tran.ID = id.Uint64()

	event := &TransactionCreated{
		CreditAmount: creditAmount,
		DebitAmount:  debitAmount,
		OrderType:    tran.OrderType,
	}

	if len(orderID) > 0 {
		tran.OrderID = orderID[0]
		event.OrderID = tran.OrderID
	}

	// fill base event data
	event.SetAggregateID(tran.ID)
	event.SetParentID(tran.UserID)
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
