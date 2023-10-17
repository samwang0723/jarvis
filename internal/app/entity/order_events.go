package entity

import "github.com/samwang0723/jarvis/internal/eventsourcing"

type OrderCreated struct {
	OrderType       string
	StockID         string
	ExchangeDate    string
	TradePrice      float32
	Quantity        uint64
	Description     string
	ProfitablePrice float32
	ProfitLoss      float32

	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*OrderCreated) EventType() eventsourcing.EventType {
	return "order.created"
}

type OrderChanged struct {
	OrderType    string
	StockID      string
	ExchangeDate string
	TradePrice   float32
	Quantity     uint64
	Description  string
	ProfitLoss   float32

	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*OrderChanged) EventType() eventsourcing.EventType {
	return "order.completed"
}

type OrderClosed struct {
	eventsourcing.BaseEvent
}

// EventType returns the name of event
func (*OrderClosed) EventType() eventsourcing.EventType {
	return "order.closed"
}
