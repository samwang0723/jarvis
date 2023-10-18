package entity

import (
	"fmt"
	"time"

	"github.com/samwang0723/jarvis/internal/eventsourcing"
)

// Define state machine
const (
	orderInitState    eventsourcing.State = ""
	orderCreatedState eventsourcing.State = "created"
	orderChangedState eventsourcing.State = "changed"
	orderClosedState  eventsourcing.State = "closed"

	taiwanStockQuantity = 1000
	dayTradeTaxRate     = 0.5
	taxRate             = 0.003
	feeRate             = 0.001425
	brokerFeeDiscount   = 0.25
	buySellTime         = 2
	percent             = 100
)

type Order struct {
	StockID          string  `gorm:"column:stock_id" json:"stockId"`
	UserID           uint64  `gorm:"column:user_id" json:"userId"`
	BuyPrice         float32 `gorm:"column:buy_price" json:"buyPrice"`
	BuyQuantity      uint64  `gorm:"column:buy_quantity" json:"buyQuantity"`
	BuyExchangeDate  string  `gorm:"column:buy_exchange_date" json:"buyExchangeDate"`
	SellPrice        float32 `gorm:"column:sell_price" json:"sellPrice"`
	SellQuantity     uint64  `gorm:"column:sell_quantity" json:"sellQuantity"`
	SellExchangeDate string  `gorm:"column:sell_exchange_date" json:"sellExchangeDate"`
	ProfitablePrice  float32 `gorm:"column:profitable_price" json:"profitablePrice"`
	Status           string  `gorm:"column:status" json:"status,omitempty"`

	ProfitLoss        float32
	ProfitLossPercent float32

	CreatedAt time.Time `gorm:"column:created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" mapstructure:"updated_at"`

	eventsourcing.BaseAggregate
}

// ensure Transaction implements Aggregate interface
var _ eventsourcing.Aggregate = &Order{}

func (Order) TableName() string {
	return "orders"
}

func (order *Order) EventTable() string {
	return "order_events"
}

func (order *Order) GetCurrentState() eventsourcing.State {
	return eventsourcing.State(order.Status)
}

func (order *Order) QuantityMatched() bool {
	return order.BuyQuantity == order.SellQuantity
}

func (order *Order) CalculateProfitLoss() {
	if order.QuantityMatched() {
		order.ProfitLoss = (order.SellPrice - order.BuyPrice) * float32(order.SellQuantity) * taiwanStockQuantity
		order.ProfitLossPercent = (order.SellPrice - order.BuyPrice) / order.BuyPrice
	}
}

//nolint:nestif // ignore nested if
func (order *Order) CalculateUnrealizedProfitLoss(currentPrice float32) {
	if order.BuyQuantity > order.SellQuantity {
		remainingQuantity := order.BuyQuantity - order.SellQuantity
		if order.SellQuantity > 0 {
			profit := (order.SellPrice - order.ProfitablePrice) * float32(order.SellQuantity) * taiwanStockQuantity
			order.ProfitLoss = profit + (currentPrice-order.ProfitablePrice)*float32(remainingQuantity)*taiwanStockQuantity
		} else {
			order.ProfitLoss = (currentPrice - order.ProfitablePrice) * float32(order.BuyQuantity) * taiwanStockQuantity
		}
		cost := order.ProfitablePrice * float32(order.BuyQuantity) * taiwanStockQuantity
		current := currentPrice * float32(order.BuyQuantity) * taiwanStockQuantity
		order.ProfitLossPercent = ((cost / current) - 1) * percent * -1
	} else if order.BuyQuantity < order.SellQuantity {
		remainingQuantity := order.SellQuantity - order.BuyQuantity
		if order.BuyQuantity > 0 {
			profit := (order.ProfitablePrice - order.BuyPrice) * float32(order.BuyQuantity) * taiwanStockQuantity
			order.ProfitLoss = profit + (order.ProfitablePrice-currentPrice)*float32(remainingQuantity)*taiwanStockQuantity
		} else {
			order.ProfitLoss = (order.ProfitablePrice - currentPrice) * float32(order.SellQuantity) * taiwanStockQuantity
		}
		cost := order.ProfitablePrice * float32(order.SellQuantity) * taiwanStockQuantity
		current := currentPrice * float32(order.SellQuantity) * taiwanStockQuantity
		order.ProfitLossPercent = ((cost / current) - 1) * percent
	}
}

// Apply updates the aggregate according to a event.
//
//nolint:lll // ignore long line length and magic number
func (order *Order) Apply(event eventsourcing.Event) error {
	newState, err := eventsourcing.TransistOnEvent(order, event)
	if err != nil {
		return fmt.Errorf("failed to transition state: %w", err)
	}

	order.Status = string(newState)

	switch event := event.(type) {
	case *OrderCreated:
		order.UserID = event.GetParentID()
		order.StockID = event.StockID
		feeAmount := event.TradePrice * float32(event.Quantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount * buySellTime
		taxAmount := event.TradePrice * float32(event.Quantity) * taiwanStockQuantity * taxRate
		originalAmount := event.TradePrice * float32(event.Quantity) * taiwanStockQuantity

		if event.OrderType == OrderTypeBuy {
			order.BuyPrice = event.TradePrice
			order.BuyQuantity = event.Quantity
			order.BuyExchangeDate = event.ExchangeDate
			order.ProfitablePrice = ((originalAmount + feeAmount + taxAmount) / float32(event.Quantity)) / taiwanStockQuantity
		} else {
			order.SellPrice = event.TradePrice
			order.SellQuantity = event.Quantity
			order.SellExchangeDate = event.ExchangeDate
			order.ProfitablePrice = ((originalAmount - feeAmount - taxAmount) / float32(event.Quantity)) / taiwanStockQuantity
		}

		order.CreatedAt = event.CreatedAt
		order.UpdatedAt = event.CreatedAt
		order.SetAggregateID(event.AggregateID)
	case *OrderChanged:
		if event.OrderType == OrderTypeBuy {
			order.BuyPrice = event.TradePrice
			order.BuyQuantity = event.Quantity
			order.BuyExchangeDate = event.ExchangeDate
		} else {
			order.SellPrice = event.TradePrice
			order.SellQuantity = event.Quantity
			order.SellExchangeDate = event.ExchangeDate
		}
		order.UpdatedAt = event.CreatedAt
	case *OrderClosed:
		order.UpdatedAt = event.CreatedAt
	default:
		return &UnsupportedEventError{event: event}
	}

	order.Version = event.GetVersion()

	return nil
}

// GetStates returns all possible state transitions
func (order *Order) GetTransitions() []eventsourcing.Transition {
	return []eventsourcing.Transition{
		{
			FromState: orderInitState,
			Event:     &OrderCreated{},
			ToState:   orderCreatedState,
		},
		{
			FromState: orderCreatedState,
			Event:     &OrderChanged{},
			ToState:   orderChangedState,
		},
		{
			FromState: orderChangedState,
			Event:     &OrderChanged{},
			ToState:   orderChangedState,
		},

		{
			FromState: orderChangedState,
			Event:     &OrderClosed{},
			ToState:   orderClosedState,
		},
	}
}

func NewOrder(
	userID uint64,
	orderType string,
	stockID string,
	exchangeDate string,
	tradePrice float32,
	quantity uint64,
) (*Order, error) {
	id, err := GenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}

	order := &Order{}
	event := &OrderCreated{
		OrderType:    orderType,
		StockID:      stockID,
		ExchangeDate: exchangeDate,
		TradePrice:   tradePrice,
		Quantity:     quantity,
	}

	// fill base event data
	event.SetAggregateID(id.Uint64())
	event.SetParentID(userID)
	event.SetVersion(1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := order.Apply(event); err != nil {
		return nil, err
	}
	// record uncommitted events
	order.AppendChange(event)

	return order, nil
}

func (order *Order) Change(
	orderType string,
	stockID string,
	exchangeDate string,
	tradePrice float32,
	quantity uint64,
) error {
	event := &OrderChanged{
		OrderType:    orderType,
		StockID:      stockID,
		ExchangeDate: exchangeDate,
		TradePrice:   tradePrice,
		Quantity:     quantity,
	}

	// fill base event data
	event.SetAggregateID(order.ID)
	event.SetParentID(order.UserID)
	event.SetVersion(order.Version + 1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := order.Apply(event); err != nil {
		return err
	}
	// record uncommitted events
	order.AppendChange(event)

	// close order if all open positions are closed
	if order.QuantityMatched() {
		order.Close()
	}

	return nil
}

func (order *Order) Close() error {
	event := &OrderClosed{}

	// fill base event data
	event.SetAggregateID(order.ID)
	event.SetParentID(order.UserID)
	event.SetVersion(order.Version + 1)
	event.SetCreatedAt(time.Now())

	// apply the event
	if err := order.Apply(event); err != nil {
		return err
	}
	// record uncommitted events
	order.AppendChange(event)

	return nil
}
