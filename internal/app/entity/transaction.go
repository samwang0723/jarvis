package entity

import (
	"fmt"
	"time"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

const (
	OrderTypeBid      = "Bid"
	OrderTypeAsk      = "Ask"
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
	StockID              string  `gorm:"column:stock_id" json:"stockId"`
	UserID               uint64  `gorm:"column:user_id" json:"userId"`
	OrderType            string  `gorm:"column:order_type" json:"orderType"`
	TradePrice           float32 `gorm:"column:trade_price" json:"tradePrice"`
	Quantity             uint64  `gorm:"column:quantity" json:"quantity"`
	ExchangeDate         string  `gorm:"column:exchange_date" json:"exchangeDate"`
	CreditAmount         float32 `gorm:"column:credit_amount" json:"creditAmount"`
	DebitAmount          float32 `gorm:"column:debit_amount" json:"debitAmount"`
	Description          string  `gorm:"column:description" json:"description"`
	ReferenceID          *uint64 `gorm:"column:reference_id" json:"referenceId"`
	Status               string  `gorm:"column:status" json:"status,omitempty"`
	OriginalExchangeDate string
	CreatedAt            time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at" mapstructure:"updated_at"`

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
	stockID string,
	userID uint64,
	orderType string,
	tradePrice float32,
	quantity uint64,
	exchangeDate string,
	creditAmount float32,
	debitAmount float32,
	description string,
	referenceID *uint64,
) (*Transaction, error) {
	id, err := GenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}

	tran := &Transaction{
		StockID:      stockID,
		UserID:       userID,
		OrderType:    orderType,
		TradePrice:   tradePrice,
		Quantity:     quantity,
		ExchangeDate: exchangeDate,
		Description:  description,
		ReferenceID:  referenceID,
	}

	tran.ID = id.Uint64()

	event := &TransactionCreated{
		CreditAmount: creditAmount,
		DebitAmount:  debitAmount,
	}

	// fill base event data
	event.SetAggregateID(tran.ID)
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
