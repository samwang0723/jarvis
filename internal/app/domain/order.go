package domain

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/helper"
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

type stockState struct {
	totalSpent    float32
	totalReceived float32
	totalFees     float32
	totalTaxes    float32
}

type Order struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	StockName        string
	StockID          string
	BuyExchangeDate  string
	Status           string
	SellExchangeDate string
	eventsourcing.BaseAggregate
	SellQuantity      uint64
	BuyQuantity       uint64
	UserID            uuid.UUID
	ProfitablePrice   float32
	SellPrice         float32
	ProfitLoss        float32
	ProfitLossPercent float32
	CurrentPrice      float32
	BuyPrice          float32
}

type ListOrdersParams struct {
	UserID        uuid.UUID
	Limit         int32
	Offset        int32
	StockIDs      []string
	Status        string
	ExchangeMonth string
}

// ensure Transaction implements Aggregate interface
var _ eventsourcing.Aggregate = &Order{}

func (order *Order) EventTable() string {
	return "order_events"
}

func (order *Order) GetCurrentState() eventsourcing.State {
	return eventsourcing.State(order.Status)
}

func (order *Order) QuantityMatched() bool {
	return order.BuyQuantity == order.SellQuantity
}

func (s *stockState) Buy(price float32, quantity uint64) {
	totalCost := price * float32(quantity) * taiwanStockQuantity
	fee := totalCost * feeRate * brokerFeeDiscount

	s.totalSpent += totalCost + fee
	s.totalFees += fee
}

func (s *stockState) Sell(price float32, quantity uint64, dayTrade bool) {
	totalRevenue := price * float32(quantity) * taiwanStockQuantity
	fee := totalRevenue * feeRate * brokerFeeDiscount
	tax := totalRevenue * taxRate
	if dayTrade {
		tax *= dayTradeTaxRate
	}

	s.totalReceived += totalRevenue - fee - tax
	s.totalFees += fee
	s.totalTaxes += tax
}

func (s *stockState) ProfitLoss() float32 {
	return helper.RoundDecimal(s.totalReceived - s.totalSpent)
}

func (s *stockState) ProfitLossPercent() float32 {
	return helper.RoundDecimalTwo((s.totalReceived - s.totalSpent) / s.totalSpent * percent)
}

func (order *Order) CalculateProfitLoss() {
	if !order.QuantityMatched() {
		return
	}

	stock := &stockState{}
	stock.Buy(order.BuyPrice, order.BuyQuantity)
	dayTrade := false
	if order.BuyExchangeDate == order.SellExchangeDate {
		dayTrade = true
	}
	stock.Sell(order.SellPrice, order.SellQuantity, dayTrade)

	order.ProfitLoss = stock.ProfitLoss()
	order.ProfitLossPercent = stock.ProfitLossPercent()
}

//nolint:nestif // ignore nested if
func (order *Order) CalculateUnrealizedProfitLoss(currentPrice float32) {
	if order.QuantityMatched() {
		return
	}

	stock := &stockState{}
	if order.BuyQuantity > order.SellQuantity {
		remainingQuantity := order.BuyQuantity - order.SellQuantity
		stock.Buy(order.BuyPrice, order.BuyQuantity)
		if order.SellQuantity > 0 {
			stock.Sell(order.SellPrice, order.SellQuantity, false)
			stock.Sell(currentPrice, remainingQuantity, false)
		} else {
			stock.Sell(currentPrice, order.BuyQuantity, false)
		}
	} else if order.BuyQuantity < order.SellQuantity {
		remainingQuantity := order.SellQuantity - order.BuyQuantity
		stock.Sell(order.SellPrice, order.SellQuantity, false)
		if order.BuyQuantity > 0 {
			stock.Buy(order.BuyPrice, order.BuyQuantity)
			stock.Buy(currentPrice, remainingQuantity)
		} else {
			stock.Buy(currentPrice, order.SellQuantity)
		}
	}

	order.ProfitLoss = stock.ProfitLoss()
	order.ProfitLossPercent = stock.ProfitLossPercent()
	order.CurrentPrice = currentPrice
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
	userID uuid.UUID,
	orderType string,
	stockID string,
	exchangeDate string,
	tradePrice float32,
	quantity uint64,
) (*Order, error) {
	id := uuid.Must(uuid.NewV4())
	order := &Order{
		BaseAggregate: eventsourcing.BaseAggregate{
			ID: id,
		},
	}
	event := &OrderCreated{
		OrderType:    orderType,
		StockID:      stockID,
		ExchangeDate: exchangeDate,
		TradePrice:   tradePrice,
		Quantity:     quantity,
	}

	// fill base event data
	event.SetAggregateID(id)
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
